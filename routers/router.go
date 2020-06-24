package routers

import (
	"Parser-Golang/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")
	beego.Router("/search/comparison", &controllers.ComparisonController{}, "post:FindComparison")
	beego.Router("/comparisons/:comp", &controllers.MainController{}, "get:DetailsPage")
	beego.Router("/update/comment/point", &controllers.ComparisonController{}, "get:UpdatePoint")

	beego.Router("/api/comparisons/:comp", &controllers.RestApiController{}, "get:GetComparisonStack")
	beego.Router("/api/livesearch/:StackName", &controllers.RestApiController{}, "get:LiveSearch")
}
