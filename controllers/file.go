package controllers

import (
	"mime"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/LyricTian/snail/utils/captcha"
	"github.com/astaxie/beego"
)

// FileController 文件下载控制器
type FileController struct {
	DownloadController
}

// Post /file
// @router /download/file [post]
func (fc *FileController) Post() {
	var reqData struct {
		FileLink   string `valid:"Required"`
		FileName   string
		FCaptchaID string `valid:"Required"`
		FCaptcha   string `valid:"Required"`
	}
	if err := fc.BindVForm(&reqData); err != nil {
		fc.Error400(err.Error())
	}
	if !captcha.VerifyString(reqData.FCaptchaID, reqData.FCaptcha) {
		fc.Error400("无效的验证码")
	}

	// 解析下载链接
	u, err := url.ParseRequestURI(reqData.FileLink)
	if err != nil {
		beego.Warn("解析下载链接错误：", err.Error())
		fc.Error400("无效的下载链接")
	}

	resp, err := fc.DoRequest(u)
	if err != nil {
		beego.Error("下载文件错误:", err.Error())
		fc.Error500("下载文件错误")
	}

	fileSize := resp.ContentLength
	fileSizeLimit := beego.AppConfig.DefaultInt64("FileSizeLimit", 512) * 1024 * 1024
	fileLink := u.String()
	fileName := fc.fileName(reqData.FileName, fileLink, resp)

	fc.AddHistory(0, int(fileSize), fileName, fileLink)

	if fileSize > 0 {
		if fileSize > fileSizeLimit {
			fc.Error400("文件大小超出限制")
		}
	}

	fc.DoResponse(fileName, fileSizeLimit, resp, nil)
}

// 获取下载文件名
func (fc *FileController) fileName(reqName, reqURL string, resp *http.Response) (filename string) {
	if reqName != "" {
		filename = reqName
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
	filename = filepath.Base(reqURL)

	return
}
