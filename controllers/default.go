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
	"strings"
	"time"

	"github.com/LyricTian/snail/models"
	"github.com/astaxie/beego"
	"github.com/dchest/captcha"
	"github.com/satori/go.uuid"
	"golang.org/x/net/http2"
)

// MainController 入口控制器
type MainController struct {
	BaseController
}

// Get /
// @router / [get]
func (mc *MainController) Get() {
	mc.Data["CaptchaID"] = captcha.NewLen(4)
	mc.LayoutSections = map[string]string{
		"Scripts": "scripts/login.html",
	}
	mc.HTML("index.html")
}

// Captcha /captcha
// @router /captcha/:id.png [get]
func (mc *MainController) Captcha() {
	id := mc.GetString(":id")
	if id == "" {
		mc.Error400("无效的验证码")
	}
	if mc.GetString("reload") != "" {
		captcha.Reload(id)
	}
	w := mc.Ctx.ResponseWriter
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	err := captcha.WriteImage(w, id, 130, 65)
	if err != nil {
		beego.Error("验证码错误：", err.Error())
		mc.Error400("无效的验证码")
	}
}

// Download /download
// @router /download [post]
func (mc *MainController) Download() {
	var reqData struct {
		FileLink  string `valid:"Required"`
		FileName  string
		CaptchaID string `valid:"Required"`
		Captcha   string `valid:"Required"`
	}
	if err := mc.BindVForm(&reqData); err != nil {
		mc.Error400(err.Error())
	}
	if !captcha.VerifyString(reqData.CaptchaID, reqData.Captcha) {
		mc.Error400("无效的验证码")
	}
	u, err := mc.parseURL(reqData.FileLink)
	if err != nil {
		mc.Error400("无效的下载链接")
		beego.Warn("解析下载链接错误：", err.Error())
	}
	resp, err := mc.doRequest(u)
	if err != nil {
		mc.Error500("下载文件错误")
		beego.Error("下载文件错误:", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		mc.Error500("请求对方服务器异常：", resp.Status)
		return
	}

	filename := reqData.FileName
	if filename = filepath.Base(u.String()); filename == "" {
		filename = mc.getRespFileName(resp)
	}

	// 增加下载历史
	history := &models.History{
		ID:       uuid.NewV4().String(),
		AccessIP: mc.Ctx.Input.IP(),
		FileLink: reqData.FileLink,
		FileName: filename,
		Time:     time.Now().Format("20060102150405"),
	}
	err = history.Create()
	if err != nil {
		beego.Info(fmt.Sprintf("%#v", history))
		beego.Error("增加下载历史发生错误：", err)
	}

	output := mc.Ctx.Output
	if l := resp.ContentLength; l > 0 {
		output.Header("Content-Length", strconv.Itoa(int(l)))
	}
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		output.Header("Content-Type", ct)
	}
	output.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(filename))
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "0")
	output.Header("Cache-Control", "must-revalidate")
	output.Header("Pragma", "public")
	io.Copy(mc.Ctx.ResponseWriter, resp.Body)
	mc.StopRun()
}

// 解析url
func (mc *MainController) parseURL(uri string) (u *url.URL, err error) {
	if !strings.Contains(uri, "://") && !strings.HasPrefix(uri, "//") {
		uri = "//" + uri
	}

	u, err = url.Parse(uri)
	if err != nil {
		return
	}
	if u.Scheme == "" {
		u.Scheme = "http"
		if !strings.HasSuffix(u.Host, ":80") {
			u.Scheme += "s"
		}
	}
	return
}

// 请求数据
func (mc *MainController) doRequest(u *url.URL) (resp *http.Response, err error) {
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
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
	return
}

// 获取响应的文件名
func (mc *MainController) getRespFileName(resp *http.Response) (filename string) {
	if hdr := resp.Header.Get("Content-Disposition"); hdr != "" {
		mt, params, err := mime.ParseMediaType(hdr)
		if err == nil && mt == "attachment" {
			if name := params["filename"]; name != "" {
				filename = name
			}
		}
	}
	if filename == "" {
		filename = uuid.NewV4().String()
	}
	return
}
