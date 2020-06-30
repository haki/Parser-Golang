package controllers

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"Parser-Golang/services"
	"github.com/astaxie/beego"
	"strings"
)

type RestApiController struct {
	beego.Controller
}

func (c *RestApiController) GetComparisonStack() {
	var comp = c.Ctx.Input.Param(":comp")

	var find bool
	comp, find = services.Parser(comp)
	if find {
		c.Data["json"] = services.ParseFromDatabase(comp)
		c.ServeJSON()
	} else {
		c.Data["json"] = "Can't Find!"
		c.ServeJSON()
	}
}

func (c *RestApiController) LiveSearch() {
	var AllStacks []models.Stack
	db.Conn.Find(&AllStacks)

	var ResultStacks []services.Stack

	StackName := c.Ctx.Input.Param(":StackName")
	StackName = strings.ToLower(StackName)
	k := 0
	for i := 0; i < len(AllStacks); i++ {
		if len(StackName) >= 1 && strings.Index(strings.ToLower(AllStacks[i].Name), StackName) != -1 && k != 15 {
			stack := services.Stack{
				Name:  AllStacks[i].Name,
				Slug:  AllStacks[i].Slug,
				Image: AllStacks[i].Image,
			}

			ResultStacks = append(ResultStacks, stack)
			k++
		}
	}

	c.Data["json"] = ResultStacks
	c.ServeJSON()
}
