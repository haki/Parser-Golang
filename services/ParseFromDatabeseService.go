package services

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
)

type Github struct {
	Star  string `json:"star"`
	Forks string `json:"forks"`
	Watch string `json:"watch"`
}

type Stats struct {
	Github Github `json:"github"`
}

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
	Stats       Stats     `json:"stats"`
	Pros        []Pros    `json:"pros"`
	Const       []Const   `json:"const"`
	Companies   []Company `json:"companies"`
}

type Comparison struct {
	Source     string  `json:"source"`
	SourcePage string  `json:"source_page"`
	Name       string  `json:"name"`
	Stacks     []Stack `json:"stacks"`
}

func ParseFromDatabase(comp string) *Comparison {
	var comparison models.Comparison
	db.Conn.Where(models.Comparison{Slug: comp}).First(&comparison)

	var stacks []models.Stack
	db.Conn.Model(&comparison).Association("Stacks").Find(&stacks)

	jsonData := &Comparison{
		Source:     "localhost:8080",
		SourcePage: "http://localhost:8080/" + comp,
		Name:       comp,
		Stacks:     []Stack{},
	}

	for i := 0; i < len(stacks); i++ {
		stack := &Stack{
			Name:        stacks[i].Name,
			Slug:        stacks[i].Slug,
			Image:       stacks[i].Image,
			Description: stacks[i].Description,
			Stats: Stats{
				Github: Github{
					Star:  stacks[i].Star,
					Forks: stacks[i].Fork,
					Watch: stacks[i].Watch,
				},
			},
			Pros:      []Pros{},
			Const:     []Const{},
			Companies: []Company{},
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

		AddStack(stack, jsonData)
	}

	return jsonData
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
	ItemStack := Stack{
		Name:        stack.Name,
		Slug:        stack.Slug,
		Image:       stack.Image,
		Description: stack.Description,
		Stats:       stack.Stats,
		Pros:        stack.Pros,
		Const:       stack.Const,
		Companies:   stack.Companies,
	}
	comparison.Stacks = append(comparison.Stacks, ItemStack)

	return comparison.Stacks
}
