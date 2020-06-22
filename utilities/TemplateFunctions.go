package utilities

import "strings"

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
