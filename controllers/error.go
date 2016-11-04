package controllers

// ErrorController 错误处理控制器
type ErrorController struct {
	BaseController
}

// Error400 400
func (ec *ErrorController) Error400() {
	ec.HTML("error/400.html")
}

// Error404 404
func (ec *ErrorController) Error404() {
	ec.HTML("error/404.html")
}

// Error500 500
func (ec *ErrorController) Error500() {
	ec.HTML("error/500.html")
}
