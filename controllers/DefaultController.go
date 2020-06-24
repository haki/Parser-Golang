package controllers

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"Parser-Golang/services"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	var allComparisons []models.Comparison
	db.Conn.Preload("Stacks").Order("view desc").Limit(15).Find(&allComparisons)

	var lastCreatedComparisons []models.Comparison
	db.Conn.Preload("Stacks").Order("id desc").Limit(15).Find(&lastCreatedComparisons)

	c.Data["TopComparisons"] = allComparisons
	c.Data["NewComparisons"] = lastCreatedComparisons

	c.TplName = "index.gohtml"
}

func (c *MainController) DetailsPage() {
	comparisonSlug := c.Ctx.Input.Param(":comp")

	find := false
	find, comparisonSlug = services.CheckComparison(comparisonSlug)

	if find {
		var comparison models.Comparison
		db.Conn.Preload("Stacks").Where(&models.Comparison{Slug: comparisonSlug}).First(&comparison)
		UpdateView(comparisonSlug)

		c.Data["Comparison"] = comparison

		c.TplName = "detailsPage.gohtml"
	} else {
		services.SaveData(comparisonSlug)
		if db.Conn.Where(&models.Comparison{Slug: comparisonSlug}).First(&models.Comparison{}).Error == nil {
			c.Redirect("/comparisons/"+comparisonSlug, 303)
		} else {
			c.Redirect("/", 303)
		}
	}
}

func UpdateView(slug string) {
	var comparison models.Comparison
	db.Conn.Where(&models.Comparison{Slug: slug}).First(&comparison)

	comparison.View = comparison.View + 1
	db.Conn.Save(&comparison)
}
