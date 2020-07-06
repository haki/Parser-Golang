package services

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func AddNewComparisons() {
	logs.Info("Update Started.")
	url := []string{"https://stackshare.io/stackups/trending", "https://stackshare.io/stackups/top", "https://stackshare.io/stackups/new"}
	for i := 0; i < len(url); i++ {
		time.Sleep(3660 * time.Millisecond)
		response, _ := http.Get(url[i])
		time.Sleep(3600 * time.Millisecond)

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
		response.Body.Close()
	}

	var AllStacks []models.Stack
	db.Conn.Find(&AllStacks)
	for i := 0; i < len(AllStacks); i++ {

		logs.Info("Getting to " + AllStacks[i].Slug)
		response, _ := http.Get("https://stackshare.io/" + AllStacks[i].Slug)
		time.Sleep(3600 * time.Millisecond)

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

		response.Body.Close()
	}

	UpdateGitData()

	logs.Info("Update successfully completed! All comparisons is up to date.")
}

func UpdateGitData() {
	logs.Info("Git Update Started.")
	var stacks []models.Stack
	db.Conn.Find(&stacks)

	for i := 0; i < len(stacks); i++ {
		if strings.Index(stacks[i].GitUrl, "github") != -1 {
			time.Sleep(3 * time.Second)
			UpdateGithubData(&stacks[i])
		}
	}

	logs.Info("Update successfully completed! Git up to date.")
}

func UpdateGithubData(stack *models.Stack) {
	type GitHub struct {
		ForksCount       int `json:"forks_count"`
		SubscribersCount int `json:"subscribers_count"`
		StargazersCount  int `json:"stargazers_count"`
	}

	gitUrl := strings.Replace(stack.GitUrl, "https://github.com/", "https://api.github.com/repos/", -1)

	response, _ := http.Get(gitUrl)
	gitJson, _ := ioutil.ReadAll(response.Body)

	github := GitHub{}
	json.Unmarshal(gitJson, &github)
	response.Body.Close()

	stack.Watch = github.SubscribersCount
	stack.Star = github.StargazersCount
	stack.Fork = github.ForksCount

	db.Conn.Save(&stack)
}

func UpdateView(slug string) {
	var comparison models.Comparison
	db.Conn.Where(&models.Comparison{Slug: slug}).First(&comparison)

	comparison.View = comparison.View + 1
	db.Conn.Save(&comparison)

	var stacks []models.Stack
	db.Conn.Model(&comparison).Association("Stacks").Find(&stacks)
	for i := 0; i < len(stacks); i++ {
		stacks[i].View = stacks[i].View + 1
		db.Conn.Save(&stacks[i])
	}
}
