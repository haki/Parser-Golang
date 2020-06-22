package routers

import (
	"Parser-Golang/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")

	beego.Router("/api/comparison/:comp", &controllers.RestApiController{}, "get:GetComparisonStack")
	beego.Router("/api/comparison/update-data", &controllers.RestApiController{}, "put:UpdateData")
	beego.Router("/api/livesearch/:StackName", &controllers.RestApiController{}, "get:LiveSearch")
}
