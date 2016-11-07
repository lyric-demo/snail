package controllers

import (
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/LyricTian/snail/models"
	"github.com/LyricTian/snail/utils/captcha"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
	"golang.org/x/net/http2"
)

// DownloadController 文件下载
type DownloadController struct {
	BaseController
}

// Post /download
// @router /download [post]
func (dc *DownloadController) Post() {
	var reqData struct {
		FileLink          string `valid:"Required"`
		FileName          string
		DownloadCaptchaID string `valid:"Required"`
		DownloadCaptcha   string `valid:"Required"`
	}
	if err := dc.BindVForm(&reqData); err != nil {
		dc.Error400(err.Error())
	}
	if !captcha.VerifyString(reqData.DownloadCaptchaID, reqData.DownloadCaptcha) {
		dc.Error400("无效的验证码")
	}

	// 解析下载链接
	u, err := url.ParseRequestURI(reqData.FileLink)
	if err != nil {
		beego.Warn("解析下载链接错误：", err.Error())
		dc.Error400("无效的下载链接")
	}

	// 请求文件数据
	resp, err := dc.doRequest(u)
	if err != nil {
		beego.Error("下载文件错误:", err.Error())
		dc.Error500("下载文件错误")
	}
	defer resp.Body.Close()

	filename := dc.getFileName(reqData.FileName, u.String(), resp)
	dc.addHistory(filename, reqData.FileLink)

	output := dc.Ctx.Output
	if l := resp.ContentLength; l > 0 {
		sizeLimit := beego.AppConfig.DefaultInt64("FileSizeLimit", 512)
		if l > sizeLimit*1024*1024 {
			dc.Error400("文件大小超出限制")
		}
		output.Header("Content-Length", strconv.Itoa(int(l)))
	}

	contentType := "application/octet-stream"
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		contentType = ct
	}
	output.Header("Content-Type", contentType)
	output.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(filename))
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "0")
	output.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	output.Header("Pragma", "no-cache")
	io.Copy(dc.Ctx.ResponseWriter, resp.Body)
	dc.StopRun()
}

// 请求数据
func (dc *DownloadController) doRequest(u *url.URL) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return
	}

	tr := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if u.Scheme == "https" {
		host, _, verr := net.SplitHostPort(req.Host)
		if verr != nil {
			host = req.Host
		}

		tr.TLSClientConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: true,
		}
		verr = http2.ConfigureTransport(tr)
		if verr != nil {
			err = verr
			return
		}
	}

	client := &http.Client{
		Transport: tr,
	}

	resp, err = client.Do(req)
	return
}

// 获取下载文件名
func (dc *DownloadController) getFileName(rfilename, rurl string, resp *http.Response) (filename string) {
	if rfilename != "" {
		filename = rfilename
		return
	}
	if hdr := resp.Header.Get("Content-Disposition"); hdr != "" {
		mt, params, err := mime.ParseMediaType(hdr)
		if err == nil && mt == "attachment" {
			if name := params["filename"]; name != "" {
				filename = name
				return
			}
		}
	}
	filename = filepath.Base(rurl)

	return
}

// 增加下载历史
func (dc *DownloadController) addHistory(filename, filelink string) {
	history := &models.History{
		ID:         uuid.NewV4().String(),
		AccessIP:   dc.Ctx.Input.IP(),
		FileLink:   filelink,
		FileName:   filename,
		CreateTime: time.Now().Format("20060102150405"),
	}
	err := history.Create()
	if err != nil {
		beego.Info(fmt.Sprintf("%#v", history))
		beego.Error("增加下载历史发生错误：", err)
	}
}
