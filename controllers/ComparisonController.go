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

	createdSlug = strings.ToLower(strings.Replace(createdSlug, " ", "-", -1))
	var find bool
	createdSlug, find = services.Parser(createdSlug)
	if find {
		c.Redirect("/comparisons/"+createdSlug, 303)
	} else {
		c.Redirect("/", 303)
	}
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

func (c *ComparisonController) AddNewComment() {
	type FormComment struct {
		Comment   string `form:"comment"`
		State     string `form:"state"`
		StackSlug string `form:"stackSlug"`
	}

	var formComment FormComment

	if err := c.ParseForm(&formComment); err == nil {
		var stack models.Stack
		db.Conn.Where(&models.Stack{Slug: formComment.StackSlug}).First(&stack)
		if formComment.State == "pros" {
			newComment := models.Pros{
				Text:    formComment.Comment,
				Point:   0,
				Enabled: false,
			}

			db.Conn.Model(&stack).Association("Pros").Append(&newComment)

		} else if formComment.State == "cons" {
			newComment := models.Cons{
				Text:    formComment.Comment,
				Point:   0,
				Enabled: false,
			}

			db.Conn.Model(&stack).Association("Cons").Append(&newComment)
		}
	}

	c.Data["json"] = "OK"
	c.ServeJSON()
}
