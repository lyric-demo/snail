package controllers

import (
	"github.com/astaxie/beego"
	"github.com/dchest/captcha"
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
		"Scripts": "scripts/index.html",
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
		beego.Warn("验证码错误：", err.Error())
		mc.Error400("验证码失效，请重新加载验证码")
	}
}
