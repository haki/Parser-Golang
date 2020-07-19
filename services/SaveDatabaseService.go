package services

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

func SaveData(comp string) (bool, string) {
	response, _ := http.Get("https://stackshare.io/stackups/" + comp)
	document, _ := goquery.NewDocumentFromReader(response.Body)
	defer response.Body.Close()

	names := strings.Split(document.Find("title").Text(), " | ")
	name := names[0]

	sourcePage, _ := document.Find("link").Attr("href")
	slug := strings.Replace(sourcePage, "https://stackshare.io/stackups/", "", -1)

	var comparison models.Comparison
	if !(strings.Index(sourcePage, "trending") != -1) && db.Conn.Where(&models.Comparison{Slug: slug}).First(&comparison).Error != nil {
		comparison = models.Comparison{
			Name:       name,
			Slug:       slug,
			View:       0,
			SourcePage: sourcePage,
		}
		db.Conn.Create(&comparison)

		SetStack(comparison, document)

		return true, comparison.Slug
	}

	return false, comp
}

func SetStack(comparison models.Comparison, document *goquery.Document) {
	db.Conn.Where(&models.Comparison{Slug: comparison.Slug}).First(&comparison)

	slugs := strings.Split(comparison.Slug, "-vs-")
	for i := 0; i < len(slugs); i++ {
		var stack models.Stack
		if db.Conn.Where(&models.Stack{Slug: slugs[i]}).First(&stack).Error != nil {
			stackResponse, _ := http.Get("https://stackshare.io/" + slugs[i])
			stackDocument, _ := goquery.NewDocumentFromReader(stackResponse.Body)

			stackName := GetName(document, stackDocument, slugs[i])

			stack = models.Stack{
				Name:        stackName,
				Slug:        slugs[i],
				Description: GetDescription(document, stackDocument, stackName),
				Image:       GetImage(document, stackDocument, slugs[i]),
				Website:     GetWebsite(stackDocument),
				GitUrl:      GetGitUrl(document, stackDocument, slugs[i]),
				Fork:        0,
				Star:        0,
				Watch:       0,
				Comparisons: nil,
				Companies:   nil,
				Cons:        nil,
				Pros:        nil,
			}

			db.Conn.Model(&comparison).Association("Stacks").Append(&stack)

			if strings.Index(stack.GitUrl, "github") != -1 {
				UpdateGithubData(&stack)
			}

			SetPros(document, stack.Name, stack.Slug)
			SetCons(document, stack.Name, stack.Slug)
			SetCompany(document, stack.Name, stack.Slug)

			stackResponse.Body.Close()
		} else {
			if stack.Name == "" {
				stackResponse, _ := http.Get("https://stackshare.io/" + stack.Slug)
				stackDocument, _ := goquery.NewDocumentFromReader(stackResponse.Body)
				stack.Name = GetName(document, stackDocument, stack.Slug)
				stack.Description = GetDescription(document, stackDocument, stack.Name)
				stack.Image = GetImage(document, stackDocument, stack.Slug)
				stack.Website = GetWebsite(stackDocument)
				stack.GitUrl = GetGitUrl(document, stackDocument, stack.Slug)
				db.Conn.Save(&stack)
			}

			if strings.Index(stack.GitUrl, "github") != -1 {
				UpdateGithubData(&stack)
			}

			db.Conn.Model(&comparison).Association("Stacks").Append(&stack)
		}
	}
}

func GetGitUrl(document *goquery.Document, stackDocument *goquery.Document, slug string) string {
	var gitUrl string

	document.Find(".css-1hlwa6q").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		if strings.Index(href, slug) != -1 {
			gitUrl = href
		}
	})

	if len(gitUrl) <= 0 {
		href, _ := stackDocument.Find(".17xwoxe a").Attr("href")
		gitUrl = href
	}

	return gitUrl
}

func GetWebsite(stackDocument *goquery.Document) string {
	var website string

	stackDocument.Find(".css-1pb731v a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		website = href
	})

	return website
}

func GetDescription(document *goquery.Document, stackDocument *goquery.Document, name string) string {
	var description string

	document.Find(".css-nil").Each(func(i int, selection *goquery.Selection) {
		if selection.Find(".css-i52n91").Text() == "What is "+name+"?" {
			description = selection.Find(".css-13sfqhu").Text()
		}
	})

	if len(description) <= 0 {
		description = stackDocument.Find(".css-1nbl3qb .css-nil .css-13sfqhu").First().Text()
	}

	return description
}

func GetImage(document *goquery.Document, stackDocument *goquery.Document, slug string) string {
	var image string

	document.Find(".css-1ogs1nl").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		if href == "/"+slug {
			src, _ := selection.Find("img").Attr("src")
			image = src
		}
	})

	if len(image) <= 0 {
		img, _ := stackDocument.Find(".css-1m5j888").Attr("src")
		image = img
	}

	return image
}

func GetName(document *goquery.Document, stackDocument *goquery.Document, slug string) string {
	var name string

	document.Find(".css-1ogs1nl").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		if href == "/"+slug {
			name = selection.Find("div").Text()
		}
	})

	if len(name) <= 0 {
		name = stackDocument.Find(".css-1cylxxa").Text()
	}

	return name
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

	document.Find(".css-4pt7vy").Each(func(i int, firstDiv *goquery.Selection) {
		var text = "What companies use " + name + "?"
		if firstDiv.Find("h2").Text() == text {
			firstDiv.Find(".css-7c9av6 .css-6nrkpz").Each(func(i int, selection *goquery.Selection) {
				companyName := selection.Find(".css-mta8ak .css-rsz8c").Text()
				companySlug, _ := selection.Find(".css-mta8ak").Attr("href")
				companyImage, _ := selection.Find(".css-mta8ak .css-1pwtf47 span .css-4lwqz5").Attr("src")

				var company models.Company
				if db.Conn.Where(&models.Company{Slug: companySlug}).First(&company).Error != nil {
					company := models.Company{
						Name:        companyName,
						Slug:        companySlug,
						Image:       companyImage,
						Website:     "",
						Email:       "",
						Github:      "",
						LinkedIn:    "",
						Facebook:    "",
						Description: "",
						Country:     "",
					}

					var optionalCompany models.Company
					if db.Conn.Where(&models.Company{Slug: companySlug}).First(&optionalCompany).Error == nil {
						db.Conn.Model(&stack).Association("Companies").Append(&optionalCompany)
					} else {
						db.Conn.Model(&stack).Association("Companies").Append(&company)
					}

				} else {
					db.Conn.Model(&stack).Association("Companies").Append(&company)
				}
			})
		}
	})
}
