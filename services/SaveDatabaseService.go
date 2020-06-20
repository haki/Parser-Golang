package services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"net/http"
	"parser/db"
	"parser/models"
	"strconv"
	"strings"
	"time"
)

func SaveData(comp string) {
	response, _ := http.Get("https://stackshare.io/stackups/" + comp)
	time.Sleep(3600 * time.Millisecond)

	if response.StatusCode != 200 {
		logs.Error("Status Code Error: " + response.Status)
		return
	} else {
		document, _ := goquery.NewDocumentFromReader(response.Body)
		defer response.Body.Close()

		sourcePage, _ := document.Find("link").Attr("href")
		name := document.Find("title").Text()
		nameTexts := strings.Split(name, "|")
		slug := strings.Replace(sourcePage, "https://stackshare.io/stackups/", "", -1)

		comparison := models.Comparison{
			Name:       nameTexts[0],
			Slug:       slug,
			View:       0,
			SourcePage: sourcePage,
			Stacks:     []models.Stack{},
		}

		SetStack(document, comparison)
	}
}

func SetStack(document *goquery.Document, comparison models.Comparison) {
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
				GithubUrl:   GetGithubData(i, document, "url"),
				Fork:        GetGithubData(i, document, "fork"),
				Star:        GetGithubData(i, document, "star"),
				Watch:       GetGithubData(i, document, "watch"),
				Cons:        nil,
				Pros:        nil,
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
		if strings.Index(strings.ToLower(selection.Find("h2").Text()), strings.ToLower(text)) != -1 {
			document.Find("ul.css-7c9av6 li").Each(func(n int, data *goquery.Selection) {
				name := data.Find("div.css-mta8ak a.css-rsz8c").Text()
				slug, _ := data.Find("div.css-mta8ak a.css-rsz8c").Attr("href")
				image, _ := data.Find("div.css-mta8ak a.css-1pwtf47 span img").Attr("src")
				if db.Conn.Where(&models.Company{Slug: slug}).First(&models.Company{}).Error != nil && slug != "" {
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

func GetGithubData(n int, document *goquery.Document, status string) string {
	var data string

	document.Find("div.css-3vlw85").Each(func(i int, selection *goquery.Selection) {
		dataNotes, _ := selection.Attr("data-notes")
		index := fmt.Sprintf("index %d num 3 offset 0", n)
		if dataNotes == index {
			githubUrl, find := selection.Find("a.css-1hlwa6q").Attr("href")
			if find {
				data = githubUrl
				if status != "url" {
					githubResp, _ := http.Get(data)
					githubDocument, _ := goquery.NewDocumentFromReader(githubResp.Body)
					defer githubResp.Body.Close()

					githubDocument.Find("div.pagehead.repohead.hx_repohead.readability-menu.bg-gray-light.pb-0.pt-3 div ul li").Each(func(i int, selection *goquery.Selection) {
						ariaLabel, _ := selection.Find("a.social-count").Attr("aria-label")
						if strings.Index(ariaLabel, "users are watching this repository") != -1 && status == "watch" {
							watch := selection.Find("a.social-count").Text()
							watchData := strings.Fields(watch)
							data = watchData[0]
						} else if strings.Index(ariaLabel, "users starred this repository") != -1 && status == "star" {
							star := selection.Find("a.social-count").Text()
							starData := strings.Fields(star)
							data = starData[0]
						} else if strings.Index(ariaLabel, "users forked this repository") != -1 && status == "fork" {
							fork := selection.Find("a.social-count").Text()
							forkData := strings.Fields(fork)
							data = forkData[0]
						}
					})
				}
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
				logs.Info(slug + " saved with success.")
			}
		})
		response.Body.Close()
	}

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

		response.Body.Close()
	}

	logs.Info("Update successfully completed! All comparisons is up to date.")
}
