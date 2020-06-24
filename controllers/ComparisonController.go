package controllers

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"Parser-Golang/services"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type ComparisonController struct {
	beego.Controller
}

func (c *ComparisonController) FindComparison() {
	type FindComparison struct {
		FirstStack  string `form:"stack0"`
		SecondStack string `form:"stack1"`
		ThirdStack  string `form:"stack2"`
	}

	var createdSlug string
	var Stacks FindComparison
	if err := c.ParseForm(&Stacks); err == nil {
		if Stacks.FirstStack != "" && Stacks.SecondStack != "" {
			if Stacks.ThirdStack != "" {
				createdSlug = Stacks.FirstStack + "-vs-" + Stacks.SecondStack + "-vs-" + Stacks.ThirdStack
			} else {
				createdSlug = Stacks.FirstStack + "-vs-" + Stacks.SecondStack
			}
		}
	}

	find := false
	find, createdSlug = services.CheckComparison(strings.ToLower(strings.Replace(createdSlug, " ", "-", -1)))
	if find {
		c.Redirect("/comparisons/"+createdSlug, 303)
	} else {
		services.SaveData(createdSlug)
		if db.Conn.Where(&models.Comparison{Slug: createdSlug}).Find(&models.Comparison{}).Error == nil {
			c.Redirect("/comparisons/"+createdSlug, 303)
		}
	}

	c.Redirect("/", 303)
}

func (c *ComparisonController) UpdatePoint() {
	var sId = c.GetString("id")
	var state = c.GetString("state")

	id, _ := strconv.Atoi(sId)

	if state == "pros" {
		var pros models.Pros
		db.Conn.Where("id = ?", id).First(&pros)
		pros.Point = pros.Point + 1
		db.Conn.Save(&pros)

	} else if state == "cons" {
		var cons models.Cons
		db.Conn.Where("id = ?", id).First(&cons)
		cons.Point = cons.Point + 1
		db.Conn.Save(&cons)
	}

	c.Data["json"] = "ok"
	c.ServeJSON()
}
