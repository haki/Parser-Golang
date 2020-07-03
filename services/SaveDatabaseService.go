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

		slugs := strings.Split(comparison.Slug, "-vs-")

		for i := 0; i < len(slugs); i++ {
			var stack models.Stack
			if db.Conn.Where(&models.Stack{Slug: slugs[i]}).First(&stack).Error != nil {
				stackResponse, _ := http.Get("https://stackshare.io/" + slugs[i])
				stackDocument, _ := goquery.NewDocumentFromReader(stackResponse.Body)

				stack = models.Stack{
					Name:        GetName(stackDocument),
					Slug:        slugs[i],
					Description: GetDescription(stackDocument),
					Image:       GetImage(stackDocument, document, comparison.Name),
					Website:     GetWebsite(stackDocument),
					GitUrl:      GetGitUrl(stackDocument),
					Fork:        "",
					Star:        "",
					Watch:       "",
					Comparisons: nil,
					Companies:   nil,
					Cons:        nil,
					Pros:        nil,
				}

				db.Conn.Model(&comparison).Association("Stacks").Append(&stack)

				SetPros(document, stack.Name, stack.Slug)
				SetCons(document, stack.Name, stack.Slug)
				SetCompany(document, stack.Name, stack.Slug)

				stackResponse.Body.Close()
			} else {
				db.Conn.Model(&comparison).Association("Stacks").Append(&stack)
			}
		}

		return true, comparison.Slug
	}

	return false, comp
}

func GetGitUrl(stackDocument *goquery.Document) string {
	href, _ := stackDocument.Find(".css-mgyi0p .css-ii8qy4 .css-12i35kv .css-1mjw833 .css-a5x1lt .css-1xqysy6 .css-17xwoxe a").Attr("href")

	return href
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

func GetDescription(stackDocument *goquery.Document) string {
	var description string
	stackDocument.Find(".css-mgyi0p .css-ii8qy4 .css-z9c3fl .css-1gs0ko2 .css-1t7lufe .css-1nbl3qb .css-nil").Each(func(i int, selection *goquery.Selection) {
		if strings.Index(selection.Find(".css-i52n91").Text(), "What is ") != -1 {
			description = selection.Find(".css-13sfqhu").First().Text()
		}
	})

	return description
}

func GetName(stackDocument *goquery.Document) string {
	return stackDocument.Find(".css-1cylxxa	").Text()
}

func GetWebsite(stackDocument *goquery.Document) string {
	website, _ := stackDocument.Find(".css-mgyi0p .css-ii8qy4 .css-12i35kv .css-1mjw833 .css-a5x1lt a").Attr("href")

	return website
}

func GetImage(stackDocument *goquery.Document, document *goquery.Document, name string) string {
	image, _ := stackDocument.Find(".css-1m5j888").Attr("src")

	if len(image) <= 22 {
		document.Find(".css-1ogs1nl").Each(func(i int, selection *goquery.Selection) {
			alt, _ := selection.Find("img").Attr("alt")
			if strings.Index(alt, name) != -1 {
				img, _ := selection.Find("img").Attr("src")
				image = img
			}
		})
	}

	return image
}
