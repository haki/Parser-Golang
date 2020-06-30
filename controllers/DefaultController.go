package controllers

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"Parser-Golang/services"
	"github.com/PuerkitoBio/goquery"
	"net/http"
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

	response, _ := http.Get("http://localhost:8080/api/comparisons/" + comparisonSlug)
	document, _ := goquery.NewDocumentFromReader(response.Body)
	if document.Text() == "\"Can't Find!\"" {
		response.Body.Close()
		c.Redirect("/", 303)
	}

	response.Body.Close()
	var find bool
	comparisonSlug, find = services.Parser(comparisonSlug)
	if find {
		c.Data["Comparison"] = services.ParseFromDatabase(comparisonSlug)
		c.TplName = "detailsPage.gohtml"
	}
}

func UpdateView(slug string) {
	var comparison models.Comparison
	db.Conn.Where(&models.Comparison{Slug: slug}).First(&comparison)

	comparison.View = comparison.View + 1
	db.Conn.Save(&comparison)
}
