package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/LyricTian/snail/utils/captcha"
	"github.com/astaxie/beego"
)

// YoutubeController youtube视频下载控制器
type YoutubeController struct {
	DownloadController
}

// Post /youtube
// @router /download/youtube [post]
func (yc *YoutubeController) Post() {
	var reqData struct {
		YFileLink  string `valid:"Required"`
		YCaptchaID string `valid:"Required"`
		YCaptcha   string `valid:"Required"`
	}

	if err := yc.BindVForm(&reqData); err != nil {
		yc.Error400(err.Error())
	}
	if !captcha.VerifyString(reqData.YCaptchaID, reqData.YCaptcha) {
		yc.Error400("无效的验证码")
	}

	videoID, err := yc.getVideoID(reqData.YFileLink)
	if err != nil || videoID == "" {
		if err != nil {
			beego.Warn("解析下载链接发生错误：", err.Error())
		}
		yc.Error400("无效的下载链接")
	}

	vis, err := yc.parseVideoInfo(videoID)
	if err != nil {
		beego.Warn("获取视频信息发生错误：", err.Error())
		yc.Error400("请求视频信息发生错误")
	}

	var cvi *YVideoInfo
	for _, vi := range vis {
		if strings.HasPrefix(vi.Type, "video/mp4") {
			cvi = vi
			break
		}
	}
	if cvi == nil {
		yc.Error400("未找到有效的视频信息")
	}

	vu, err := url.ParseRequestURI(cvi.URL)
	if err != nil {
		beego.Warn("解析视频下载链接发生错误：", err.Error())
		yc.Error400("解析视频下载链接发生错误")
	}

	resp, err := yc.DoRequest(vu)
	if err != nil {
		beego.Warn("请求视频下载发生错误：", err.Error())
		yc.Error400("请求视频下载发生错误")
	}

	fileSize := resp.ContentLength
	fileName := fmt.Sprintf("%s.mp4", cvi.Title)
	fileSizeLimit := beego.AppConfig.DefaultInt64("YoutubeLimit", 512) * 1024 * 1024

	// 增加文件下载历史
	yc.AddHistory(1, int(fileSize), fileName, reqData.YFileLink)

	if fileSize > 0 {
		if fileSize > fileSizeLimit {
			yc.Error400("视频文件大小超出限制")
		}
	}

	yc.DoResponse(fileName, fileSizeLimit, resp, nil)
}

// 解析视频信息
func (yc *YoutubeController) parseVideoInfo(videoID string) (videoInfos []*YVideoInfo, err error) {
	rawurl := fmt.Sprintf("http://www.youtube.com/get_video_info?video_id=%s", videoID)

	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return
	}
	resp, err := yc.DoRequest(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	uvs, err := url.ParseQuery(string(buf))
	if err != nil {
		return
	}

	if uvs.Get("status") == "fail" {
		err = errors.New(uvs.Get("reason"))
		return
	}

	streams := strings.Split(uvs.Get("url_encoded_fmt_stream_map"), ",")
	for _, stream := range streams {
		suvs, verr := url.ParseQuery(stream)
		if verr != nil {
			err = verr
			return
		}
		videoInfo := &YVideoInfo{
			Title:   uvs.Get("title"),
			Author:  uvs.Get("author"),
			URL:     suvs.Get("url"),
			Quality: suvs.Get("quality"),
			Type:    suvs.Get("type"),
		}
		videoInfos = append(videoInfos, videoInfo)
	}

	return
}

// 获取视频ID
func (yc *YoutubeController) getVideoID(fileLink string) (videoID string, err error) {
	u, err := url.ParseRequestURI(fileLink)
	if err != nil {
		return
	} else if !strings.Contains(u.Host, "youtube.com") {
		return
	}
	videoID = u.Query().Get("v")
	return
}

// YVideoInfo youtube视频信息
type YVideoInfo struct {
	Title   string
	Author  string
	URL     string
	Quality string
	Type    string
}
