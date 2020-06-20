package routers

import (
	"github.com/astaxie/beego"
	"parser/controllers"
)

func init() {
   	beego.Router("/api/comparison/:comp", &controllers.RestApiController{}, "get:GetComparisonStack")
   	beego.Router("/api/comparison/update-data", &controllers.RestApiController{}, "get:UpdateData")
}
