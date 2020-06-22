package controllers

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	var allComparisons []models.Comparison
	db.Conn.Preload("Stacks").Order("view desc").Limit(15).Find(&allComparisons)

	c.Data["TopComparisons"] = allComparisons

	c.TplName = "index.gohtml"
}
