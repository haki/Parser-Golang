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

func AddNewComparisons() {
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
				logs.Info(slug + " saved with success.")
			}
		})
		DeleteComparisonIfHasProblem()
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
				logs.Info(slug + " saved with success.")
			}
		})
		DeleteComparisonIfHasProblem()

		response.Body.Close()
	}

	logs.Info("Update successfully completed! All comparisons is up to date.")
}

func UpdateGitData() {
	var stack []models.Stack
	db.Conn.Find(&stack)

	for i := 0; i < len(stack); i++ {
		if strings.Index(stack[i].GitUrl, "github") != -1 {
			response, _ := http.Get(stack[i].GitUrl)
			document, _ := goquery.NewDocumentFromReader(response.Body)

			n := 0
			var data [3]string
			document.Find("ul.pagehead-actions li").Each(func(k int, selection *goquery.Selection) {
				fields := strings.Fields(selection.Find("a.social-count").Text())
				if len(fields) >= 1 {
					data[n] = fields[0]
					n++
				}
				if n == 3 {
					stack[i].Watch = data[0]
					stack[i].Star = data[1]
					stack[i].Fork = data[2]
					db.Conn.Save(&stack[i])
					n = 0
				}
			})

			response.Body.Close()
		}
	}

	logs.Info("Update successfully completed! Git up to date.")
}

func DeleteComparisonIfHasProblem() {
	var comparison []models.Comparison
	db.Conn.Find(&comparison)

	for i := 0; i < len(comparison); i++ {
		if db.Conn.Model(&comparison[i]).Association("Stacks").Count() <= 1 || comparison[i].SourcePage == "https://stackshare.io/stackups/trending" {
			db.Conn.Model(&comparison[i]).Association("Stacks").Delete()
			db.Conn.Unscoped().Delete(&comparison[i])
		}
	}
}
