package controllers

import (
	"mime"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/dchest/captcha"
)

// FileController 文件下载控制器
type FileController struct {
	DownloadController
}

// Post /file
// @router /download/file [post]
func (fc *FileController) Post() {
	var reqData struct {
		FileLink          string `valid:"Required"`
		FileName          string
		DownloadCaptchaID string `valid:"Required"`
		DownloadCaptcha   string `valid:"Required"`
	}
	if err := fc.BindVForm(&reqData); err != nil {
		fc.Error400(err.Error())
	}
	if !captcha.VerifyString(reqData.DownloadCaptchaID, reqData.DownloadCaptcha) {
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
	defer resp.Body.Close()

	fileSize := resp.ContentLength
	if fileSize > 0 {
		sizeLimit := beego.AppConfig.DefaultInt64("FileSizeLimit", 512)
		if fileSize > sizeLimit*1024*1024 {
			fc.Error400("文件大小超出限制")
		}
	}

	fileLink := u.String()
	fileName := fc.fileName(reqData.FileName, fileLink, resp)
	fc.AddHistory(0, int(fileSize), fileName, fileLink)

	fc.DoResponse(fileName, resp, nil)
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
