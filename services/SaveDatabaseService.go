package services

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SaveData(comp string) {
	response, err := http.Get("https://stackshare.io/stackups/" + comp)
	time.Sleep(3600 * time.Millisecond)

	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		logs.Error("Status Code Error: " + response.Status)
		return
	} else {
		document, _ := goquery.NewDocumentFromReader(response.Body)
		defer response.Body.Close()

		sourcePage, _ := document.Find("link").Attr("href")
		if sourcePage == "https://stackshare.io/stackups/trending" {
			return
		}

		name := document.Find("title").Text()
		nameTexts := strings.Split(name, "|")
		slug := strings.Replace(sourcePage, "https://stackshare.io/stackups/", "", -1)

		comparison := &models.Comparison{
			Name:       nameTexts[0],
			Slug:       slug,
			View:       0,
			SourcePage: sourcePage,
			Stacks:     []models.Stack{},
		}

		if sourcePage != "https://stackshare.io/stackups/trending" {
			SetStack(document, comparison)
		} else {
			return
		}
	}
}

func SetStack(document *goquery.Document, comparison *models.Comparison) {
	db.Conn.Create(&comparison)
	db.Conn.Where(&models.Comparison{Slug: comparison.Slug}).First(&comparison)
	document.Find("div.css-x7ngfe a.css-1ogs1nl").Each(func(i int, stackItem *goquery.Selection) {
		name := stackItem.Find("div").Text()
		slug, _ := stackItem.Attr("href")
		slug = strings.Replace(slug, "/", "", -1)
		image, _ := stackItem.Find("img").Attr("src")

		db.Conn.Where(&models.Stack{Slug: slug}).First(&models.Stack{})
		if err := db.Conn.Where(&models.Stack{Slug: slug}).First(&models.Stack{}).Error; err != nil {
			stack := models.Stack{
				Name:        name,
				Slug:        slug,
				Description: GetDescription(i, document),
				Image:       image,
				GitUrl:      GetGitUrl(i, document),
			}

			if optional := db.Conn.Where(&models.Stack{Slug: slug}).First(&models.Stack{}).Error; optional == nil {
				var optionalStack models.Stack
				db.Conn.Where(&models.Stack{Slug: slug}).First(&optionalStack)
				db.Conn.Model(&comparison).Association("Stacks").Append(&optionalStack)
			} else {
				db.Conn.Model(&comparison).Association("Stacks").Append(&stack)
			}

			SetCompany(document, name, slug)
			SetPros(document, name, slug)
			SetCons(document, name, slug)

		} else {
			var stack models.Stack
			db.Conn.Where(&models.Stack{Slug: slug}).First(&stack)
			db.Conn.Model(&comparison).Association("Stacks").Append(&stack)
		}
	})
}

func SetCons(document *goquery.Document, name string, slug string) {
	var stack models.Stack
	db.Conn.Where(&models.Stack{Slug: slug}).First(&stack)

	title := "Cons of " + name
	document.Find("div.css-3vlw85").Each(func(i int, selection *goquery.Selection) {
		if strings.Index(selection.Find("h2").Text(), title) != -1 {
			selection.Find("ul.css-7c9av6 li").Each(func(n int, readCons *goquery.Selection) {
				strPoint := readCons.Find("div div span").Text()
				strPoint = strings.Replace(strPoint, ".", "", -1)
				strPoint = strings.Replace(strPoint, "k", "", -1)
				point, _ := strconv.Atoi(strPoint)
				text := readCons.Find("div a").Text()

				cons := models.Cons{
					Text:    text,
					Point:   point,
					Enabled: true,
				}

				db.Conn.Model(&stack).Association("Cons").Append(&cons)
			})
		}
	})
}

func SetPros(document *goquery.Document, name string, slug string) {
	var stack models.Stack
	db.Conn.Where(&models.Stack{Slug: slug}).First(&stack)

	title := "Pros of " + name
	document.Find("div.css-3vlw85 div.css-nil div.css-1v4wqws div.css-nil div.css-uxqild").Each(func(i int, selection *goquery.Selection) {
		if strings.Index(selection.Find("h2").Text(), title) != -1 {
			selection.Find("ul.css-7c9av6 li").Each(func(n int, readPro *goquery.Selection) {
				strPoint := readPro.Find("span.css-5x5cr6").Text()
				strPoint = strings.Replace(strPoint, ".", "", -1)
				strPoint = strings.Replace(strPoint, "k", "", -1)
				point, _ := strconv.Atoi(strPoint)
				text := readPro.Find("a.css-1iibd1t").Text()

				pros := models.Pros{
					Text:    text,
					Point:   point,
					Enabled: true,
				}

				db.Conn.Model(&stack).Association("Pros").Append(&pros)
			})
		}
	})
}

func SetCompany(document *goquery.Document, name string, slug string) {
	var stack models.Stack
	db.Conn.Where(&models.Stack{Slug: slug}).First(&stack)

	var text = "What companies use " + name + "?"

	document.Find("div.css-3vlw85").Each(func(i int, selection *goquery.Selection) {

		if selection.Find("h2.css-nil").Text() == text {
			selection.Find("ul.css-7c9av6 li").Each(func(n int, data *goquery.Selection) {

				name := data.Find("div.css-mta8ak a.css-rsz8c").Text()
				slug, _ := data.Find("div.css-mta8ak a.css-rsz8c").Attr("href")
				image, _ := data.Find("div.css-mta8ak a.css-1pwtf47 span img").Attr("src")
				if db.Conn.Where(&models.Company{Slug: slug}).First(&models.Company{}).Error != nil {
					company := models.Company{
						Name:        name,
						Slug:        slug,
						Image:       image,
						Website:     "",
						Email:       "",
						Github:      "",
						LinkedIn:    "",
						Facebook:    "",
						Description: "",
						Country:     "",
					}

					if db.Conn.Where(&models.Company{Slug: slug}).First(&models.Company{}).Error == nil {
						var optionalCompany models.Company
						db.Conn.Where(&models.Company{Slug: slug}).First(&optionalCompany)
						db.Conn.Model(&stack).Association("Companies").Append(&optionalCompany)
					} else {
						db.Conn.Model(&stack).Association("Companies").Append(&company)
					}

				} else {
					var company models.Company
					db.Conn.Where(&models.Company{Slug: slug}).First(&company)
					db.Conn.Model(&stack).Association("Companies").Append(&company)
				}
			})
		}
	})
}

func GetGitUrl(n int, document *goquery.Document) string {
	var data string

	document.Find("div.css-3vlw85").Each(func(i int, selection *goquery.Selection) {
		dataNotes, _ := selection.Attr("data-notes")
		index := fmt.Sprintf("index %d num 3 offset 0", n)
		if dataNotes == index {
			githubUrl, find := selection.Find("a.css-1hlwa6q").Attr("href")
			if find {
				data = githubUrl
			}
		}
	})

	return data
}

func GetDescription(n int, document *goquery.Document) string {
	var description string
	document.Find("div.css-3vlw85").Each(func(i int, selection *goquery.Selection) {
		dataNotes, _ := selection.Attr("data-notes")
		index := fmt.Sprintf("index %d num 3 offset 0", n)

		if dataNotes == index && strings.Index(selection.Find("h2.css-i52n91").Text(), "What is") != -1 {
			description = selection.Find("div.css-nil div.css-13sfqhu").Text()
		}
	})

	return description
}
