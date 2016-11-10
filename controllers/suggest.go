package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LyricTian/snail/models"
	"github.com/LyricTian/snail/utils/captcha"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
)

// SuggestController 反馈建议
type SuggestController struct {
	BaseController
}

// Post /suggest
// @router /suggest [post]
func (sc *SuggestController) Post() {
	var reqData struct {
		Email      string `valid:"Required"`
		Comment    string `valid:"Required"`
		SCaptchaID string `valid:"Required"`
		SCaptcha   string `valid:"Required"`
	}
	if err := sc.BindVForm(&reqData); err != nil {
		sc.Error400(err.Error())
	}
	if !captcha.VerifyString(reqData.SCaptchaID, reqData.SCaptcha) {
		sc.Error400("无效的验证码")
	}
	suggest := &models.Suggest{
		ID:         uuid.NewV4().String(),
		AccessIP:   sc.Ctx.Input.IP(),
		Email:      reqData.Email,
		Comment:    reqData.Comment,
		CreateTime: time.Now().Format("20060102150405"),
	}
	err := suggest.Create()
	if err != nil {
		beego.Info(fmt.Sprintf("%#v", suggest))
		beego.Error("增加反馈建议发生错误：", err.Error())
	}
	sc.SetSession("suggest", "true")
	sc.Redirect("/suggest", http.StatusFound)
}

// Get /suggest
// @router /suggest [get]
func (sc *SuggestController) Get() {
	if sc.GetSession("suggest") == nil {
		sc.Redirect("/", http.StatusFound)
		return
	}
	sc.DelSession("suggest")
	sc.HTML("suggest.html")
}
