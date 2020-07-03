package controllers

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
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

	var findComparison FindComparison
	var comparisonSlug string
	if err := c.ParseForm(&findComparison); err == nil {
		comparisonSlug = findComparison.FirstStack + "-vs-" + findComparison.SecondStack
		if findComparison.ThirdStack != "" {
			comparisonSlug = comparisonSlug + "-vs-" + findComparison.ThirdStack
		}
		comparisonSlug = strings.Replace(comparisonSlug, " ", "-", -1)
		comparisonSlug = strings.Replace(comparisonSlug, ".", "", -1)
		comparisonSlug = strings.ToLower(comparisonSlug)
		c.Redirect("/comparisons/"+comparisonSlug, 303)
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

func (c *ComparisonController) AddNewComment() {
	type FormComment struct {
		ComparisonSlug string `form:"comparisonSlug"`
		Comment        string `form:"comment"`
		State          string `form:"state"`
		StackSlug      string `form:"stackSlug"`
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

	c.Redirect("/comparisons/"+formComment.ComparisonSlug, 302)
}
