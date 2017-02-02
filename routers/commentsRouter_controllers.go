package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:FileController"] = append(beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:FileController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/download/file`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:MainController"] = append(beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:MainController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:MainController"] = append(beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:MainController"],
		beego.ControllerComments{
			Method: "Captcha",
			Router: `/captcha/:id.png`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:SuggestController"] = append(beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:SuggestController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/suggest`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:SuggestController"] = append(beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:SuggestController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/suggest`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:YoutubeController"] = append(beego.GlobalControllerRouter["github.com/LyricTian/snail/controllers:YoutubeController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/download/youtube`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
