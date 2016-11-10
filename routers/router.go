package routers

import (
	"github.com/LyricTian/snail/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Include(new(controllers.MainController))
	beego.Include(new(controllers.FileController))
	beego.Include(new(controllers.YoutubeController))
	beego.Include(new(controllers.SuggestController))
}
