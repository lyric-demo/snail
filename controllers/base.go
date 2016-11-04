package controllers

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

// BaseController 基础控制器
type BaseController struct {
	beego.Controller
}

// JSON 响应json数据
func (bc *BaseController) JSON(v interface{}) {
	bc.Data["json"] = v
	bc.ServeJSON()
}

// HTML 响应html渲染模板
func (bc *BaseController) HTML(tplName string) {
	bc.Layout = "share/layout.html"
	bc.TplName = tplName
}

// Error400 响应400错误码
func (bc *BaseController) Error400(msgs ...string) {
	body := "请求发生错误"
	if len(msgs) > 0 {
		body = msgs[0]
	}
	bc.Data["Message"] = body
	bc.Abort("400")
}

// Error500 响应500错误码
func (bc *BaseController) Error500(msgs ...string) {
	body := "服务端异常"
	if len(msgs) > 0 {
		body = msgs[0]
	}
	bc.Data["Message"] = body
	bc.Abort("500")
}

// BindVForm 绑定并验证表单数据
func (bc *BaseController) BindVForm(obj interface{}) (err error) {
	if verr := bc.ParseForm(obj); verr != nil {
		err = errors.New("无效的请求参数")
		beego.Warn("请求参数错误：", verr.Error())
		return
	}
	valid := new(validation.Validation)
	if b, verr := valid.Valid(obj); verr != nil || !b {
		err = errors.New("无效的请求数据")
		if verr != nil {
			beego.Warn("验证请求数据错误：", verr.Error())
			return
		}
		if len(valid.Errors) > 0 {
			beego.Warn("验证请求数据错误：", valid.Errors[0].Field+",", valid.Errors[0].Error())
		}
	}
	return
}
