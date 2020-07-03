package services

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
)

type Pros struct {
	Text  string `json:"text"`
	Point int    `json:"point"`
}

type Const struct {
	Text  string `json:"text"`
	Point int    `json:"point"`
}

type Company struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image string `json:"image"`
}

type Stack struct {
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	GitUrl      string    `json:"git_url"`
	Star        string    `json:"star"`
	Fork        string    `json:"fork"`
	Watch       string    `json:"watch"`
	Website     string    `json:"website"`
	Pros        []Pros    `json:"pros"`
	Const       []Const   `json:"const"`
	Companies   []Company `json:"companies"`
}

type Comparison struct {
	Source     string  `json:"source"`
	SourcePage string  `json:"source_page"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	Stacks     []Stack `json:"stacks"`
}

func ParseFromDatabase(comp string) Comparison {
	var comparison models.Comparison
	db.Conn.Where(models.Comparison{Slug: comp}).First(&comparison)

	var stacks []models.Stack
	db.Conn.Model(&comparison).Association("Stacks").Find(&stacks)

	jsonPComparison := &Comparison{
		Source:     "localhost:8080",
		SourcePage: "http://localhost:8080/" + comp,
		Name:       comparison.Name,
		Slug:       comparison.Slug,
		Stacks:     []Stack{},
	}

	for i := 0; i < len(stacks); i++ {
		stack := &Stack{
			Name:        stacks[i].Name,
			Slug:        stacks[i].Slug,
			Image:       stacks[i].Image,
			Description: stacks[i].Description,
			GitUrl:      stacks[i].GitUrl,
			Star:        stacks[i].Star,
			Fork:        stacks[i].Fork,
			Watch:       stacks[i].Watch,
			Website:     stacks[i].Website,
			Pros:        []Pros{},
			Const:       []Const{},
			Companies:   []Company{},
		}

		var pros []models.Pros
		db.Conn.Model(&stacks[i]).Association("Pros").Find(&pros)

		for i := 0; i < len(pros); i++ {
			newPros := Pros{
				Text:  pros[i].Text,
				Point: pros[i].Point,
			}

			AddPros(newPros, stack)
		}

		var cons []models.Cons
		db.Conn.Model(&stacks[i]).Association("Cons").Find(&cons)

		for i := 0; i < len(cons); i++ {
			newCons := Const{
				Text:  cons[i].Text,
				Point: cons[i].Point,
			}

			AddCons(newCons, stack)
		}

		var companies []models.Company
		db.Conn.Model(&stacks[i]).Association("Companies").Find(&companies)

		for i := 0; i < len(companies); i++ {
			newCompany := Company{
				Name:  companies[i].Name,
				Slug:  companies[i].Slug,
				Image: companies[i].Image,
			}

			AddCompany(newCompany, stack)
		}

		AddStack(stack, jsonPComparison)
	}

	jsonData := *jsonPComparison

	return jsonData
}

func TopComparisons() [15]Comparison {
	var topComparisons []models.Comparison
	db.Conn.Preload("Stacks").Order("view desc").Limit(15).Find(&topComparisons)

	var jsonTopComparisons [15]Comparison

	for i := 0; i < len(topComparisons); i++ {
		jsonTopComparisons[i] = ParseFromDatabase(topComparisons[i].Slug)
	}

	return jsonTopComparisons
}

func NewComparisons() [15]Comparison {
	var newComparisons []models.Comparison
	db.Conn.Preload("Stacks").Order("id desc").Limit(15).Find(&newComparisons)

	var jsonNewComparisons [15]Comparison

	for i := 0; i < len(newComparisons); i++ {
		jsonNewComparisons[i] = ParseFromDatabase(newComparisons[i].Slug)
	}

	return jsonNewComparisons
}

func AddCompany(company Company, stack *Stack) []Company {
	stack.Companies = append(stack.Companies, company)
	return stack.Companies
}

func AddPros(stackItem Pros, stack *Stack) []Pros {
	stack.Pros = append(stack.Pros, stackItem)
	return stack.Pros
}

func AddCons(stackItem Const, stack *Stack) []Const {
	stack.Const = append(stack.Const, stackItem)
	return stack.Const
}

func AddStack(stack *Stack, comparison *Comparison) []Stack {
	item := *stack
	comparison.Stacks = append(comparison.Stacks, item)

	return comparison.Stacks
}
