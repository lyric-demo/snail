package controllers

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/LyricTian/snail/models"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
	"golang.org/x/net/http2"
)

// DownloadController 文件下载
type DownloadController struct {
	BaseController
}

// DoRequest 请求数据
func (dc *DownloadController) DoRequest(u *url.URL) (resp *http.Response, err error) {
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
		Timeout:   time.Second * 30,
	}
	resp, err = client.Do(req)
	return
}

// DoResponse 响应文件数据
func (dc *DownloadController) DoResponse(fileName string, fileLimit int64, resp *http.Response, header http.Header) {
	defer resp.Body.Close()

	output := dc.Ctx.Output
	if cl := resp.ContentLength; cl > 0 {
		output.Header("Content-Length", strconv.Itoa(int(cl)))
	}
	contentType := "application/octet-stream"
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		contentType = ct
	}
	output.Header("Content-Type", contentType)
	fileName = url.QueryEscape(fileName)
	output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s;filename*=utf-8''%s", fileName, fileName))
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "0")
	output.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	output.Header("Pragma", "no-cache")
	if header != nil {
		for key := range header {
			output.Header(key, header.Get(key))
		}
	}

	reader := &io.LimitedReader{
		R: resp.Body,
	}
	if fileLimit > 0 {
		reader.N = fileLimit
	}
	io.Copy(dc.Ctx.ResponseWriter, reader)
	dc.StopRun()
}

// AddHistory 增加下载历史
func (dc *DownloadController) AddHistory(fileType, fileSize int, fileName, fileLink string) {
	history := &models.History{
		ID:         uuid.NewV4().String(),
		AccessIP:   dc.Ctx.Input.IP(),
		FileLink:   fileLink,
		FileName:   fileName,
		FileSize:   fileSize,
		FileType:   fileType,
		CreateTime: time.Now().Format("20060102150405"),
	}
	err := history.Create()
	if err != nil {
		beego.Info(fmt.Sprintf("%#v", history))
		beego.Error("增加下载历史发生错误：", err)
	}
}
