package routers

import (
	"controllers"

	"github.com/astaxie/beego"
)

func init() {
	//设置路由
	beego.Router("/", &controllers.MainController{})
	beego.Router("/v:ver:int/*", &controllers.API1{})
	beego.Router("/:namespace/:resource/:action", &controllers.Tutorial{})
	beego.Router("/user/:userid/:function/:obj/:act", &controllers.UserInfo{})
}
