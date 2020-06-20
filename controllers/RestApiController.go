package controllers

import (
	"github.com/astaxie/beego"
	"parser/models"
	"parser/services"
	"strings"
)

type RestApiController struct {
	beego.Controller
}

var comparison models.Comparison
var stack []models.Stack

func (c *RestApiController) GetComparisonStack() {
	var comp = c.Ctx.Input.Param(":comp")

	c.Data["json"] = services.Parser(strings.ToLower(comp))
	c.ServeJSON()
}

func (c *RestApiController) UpdateData() {
	services.UpdateData()
	c.Data["json"] = "OK!"
	c.ServeJSON()
}


