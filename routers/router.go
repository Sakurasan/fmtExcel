package routers

import (
	"fmtExcel/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/query", &controllers.QueryController{})

}
