package main

import (
	"fmt"

	"github.com/LyricTian/snail/controllers"
	"github.com/LyricTian/snail/models"
	"github.com/astaxie/beego"

	_ "github.com/LyricTian/snail/routers"
)

func main() {
	// 初始化下载历史DB
	models.InitHistoryDB(beego.AppConfig.String("db::History"))
	defer models.HistoryDB.Close()

	// 设定日志
	beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s"}`, beego.AppConfig.String("LogFile")))

	// 错误处理
	beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}
