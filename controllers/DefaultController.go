package controllers

import (
	"Parser-Golang/services"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	c.Data["TopComparisons"] = services.TopComparisons()
	c.Data["NewComparisons"] = services.NewComparisons()

	c.TplName = "index.gohtml"
}

func (c *MainController) DetailsPage() {
	comparisonSlug := c.Ctx.Input.Param(":comp")

	var find bool
	comparisonSlug, find = services.Parser(comparisonSlug)

	if !find {
		c.Redirect("/", 303)
	} else {
		services.UpdateView(comparisonSlug)
		c.Data["Comparison"] = services.ParseFromDatabase(comparisonSlug)
		c.TplName = "detailsPage.gohtml"
	}
}
