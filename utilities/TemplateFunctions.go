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
