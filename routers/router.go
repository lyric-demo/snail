package routers

import (
	"github.com/LyricTian/snail/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Include(new(controllers.MainController))
}
