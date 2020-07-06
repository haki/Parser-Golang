package utilities

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"strings"
)

func IfNotZero(value int) bool {
	if value != 0 {
		return true
	}
	return false
}

func GetFirstWord(word string) string {
	wordArr := strings.Split(word, " ")

	return wordArr[0]
}

func PreloadStackPros(stackSlug string) models.Stack {
	var stack models.Stack
	db.Conn.Preload("Pros").Where(&models.Stack{Slug: stackSlug}).Find(&stack)

	return stack
}

func PreloadStackCons(stackSlug string) models.Stack {
	var stack models.Stack
	db.Conn.Preload("Cons").Where(&models.Stack{Slug: stackSlug}).Find(&stack)

	return stack
}

func PreloadStackCompanies(stackSlug string) models.Stack {
	var stack models.Stack
	db.Conn.Preload("Companies").Where(&models.Stack{Slug: stackSlug}).Find(&stack)

	return stack
}

func ChartData(comparisonSlug string) []int64 {
	var data []int64
	var comparison models.Comparison
	db.Conn.Where(&models.Comparison{Slug: comparisonSlug}).First(&comparison)

	var stacks []models.Stack
	db.Conn.Model(&comparison).Association("Stacks").Find(&stacks)

	for i := 0; i < len(stacks); i++ {
		data = append(data, stacks[i].View)
	}

	data = append(data, 0)

	return data
}

func ChartLabels(comparisonSlug string) []string {
	var data []string
	var comparison models.Comparison
	db.Conn.Where(&models.Comparison{Slug: comparisonSlug}).First(&comparison)

	var stacks []models.Stack
	db.Conn.Model(&comparison).Association("Stacks").Find(&stacks)

	for i := 0; i < len(stacks); i++ {
		data = append(data, stacks[i].Name)
	}

	return data
}
