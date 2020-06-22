package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {
	beego.BConfig.Listen.ServerTimeOut = 0
	c.Data["RunMode"] = beego.AppConfig.String("runmode")
	c.Layout = "layouts/app.gohtml"
}
