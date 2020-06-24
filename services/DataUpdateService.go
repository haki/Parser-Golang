package services

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"net/http"
	"strings"
	"time"
)

func UpdateData() {
	logs.Info("Update Started.")
	url := []string{"https://stackshare.io/stackups/trending", "https://stackshare.io/stackups/top", "https://stackshare.io/stackups/new"}
	for i := 0; i < len(url); i++ {
		time.Sleep(3600 * time.Millisecond)
		response, _ := http.Get(url[i])
		document, _ := goquery.NewDocumentFromReader(response.Body)

		document.Find("div.grid-item a").Each(func(i int, comparisons *goquery.Selection) {
			slug, _ := comparisons.Attr("href")
			slug = strings.Replace(slug, "/stackups/", "", -1)
			if db.Conn.Where(&models.Comparison{Slug: slug}).Find(&models.Comparison{}).Error != nil {
				logs.Info("Parsing From " + slug + "...")
				SaveData(slug)
				DeleteComparisonIfHasOneStack(slug)
				logs.Info(slug + " saved with success.")
			}
		})
		response.Body.Close()
	}

	time.Sleep(5 * time.Minute)

	var AllStacks []models.Stack
	db.Conn.Find(&AllStacks)
	for i := 0; i < len(AllStacks); i++ {

		logs.Info("Getting to " + AllStacks[i].Slug)
		time.Sleep(3600 * time.Millisecond)
		response, _ := http.Get("https://stackshare.io/" + AllStacks[i].Slug)
		document, _ := goquery.NewDocumentFromReader(response.Body)
		document.Find("div.css-nuwf1p div.css-13zfms0 div.css-1rmabp8 a").Each(func(k int, selection *goquery.Selection) {
			slug, _ := selection.Attr("href")
			slug = strings.Replace(slug, "/stackups/", "", -1)
			if db.Conn.Where(&models.Comparison{Slug: slug}).Find(&models.Comparison{}).Error != nil {
				logs.Info("Parsing from " + slug + "...")
				SaveData(slug)
				DeleteComparisonIfHasOneStack(slug)
				logs.Info(slug + " saved with success.")
			}
		})

		response.Body.Close()
	}

	logs.Info("Update successfully completed! All comparisons is up to date.")
}

func UpdateGitData() {

}

func DeleteComparisonIfHasOneStack(slug string) {
	var comparison models.Comparison
	if db.Conn.Where(&models.Comparison{Slug: slug}).Find(&comparison).Error == nil {
		if db.Conn.Model(&comparison).Association("Stacks").Count() <= 1 {
			db.Conn.Model(&comparison).Association("Stacks").Delete()
			db.Conn.Unscoped().Delete(&comparison)
		}
	}
}
