package services

import (
	"Parser-Golang/db"
	"Parser-Golang/models"
	"strings"
)

func Parser(comp string) *Comparison {
	var find = false
	find, comp = CheckComparison(comp)

	if !find {
		SaveData(comp)
	}

	return ParseFromDatabase(comp)
}

func CheckComparison(comp string) (bool, string) {
	find := false
	slug := strings.Split(comp, "-vs-")
	for i := 0; i < len(slug); i++ {
		var stack models.Stack
		if db.Conn.Where(&models.Stack{Slug: slug[i]}).First(&stack).Error == nil && !find {
			if len(slug) == 2 {
				find, comp = CheckComparisonWith2Stacks(slug)
			} else if len(slug) == 3 {
				find, comp = CheckComparisonWith3Stacks(slug)
			}
		} else {
			break
		}
	}

	return find, comp
}

func CheckComparisonWith2Stacks(slug []string) (bool, string) {
	var comparison models.Comparison
	for k := 0; k < len(slug); k++ {
		for l := 0; l < len(slug); l++ {
			compSlug := slug[k] + "-vs-" + slug[l]
			if db.Conn.Where(&models.Comparison{Slug: compSlug}).First(&comparison).Error == nil {
				return true, compSlug
			}
		}
	}

	compSlug := slug[0] + "-vs-" + slug[1]
	return false, compSlug
}

func CheckComparisonWith3Stacks(slug []string) (bool, string) {
	var comparison models.Comparison
	for k := 0; k < len(slug); k++ {
		for l := 0; l < len(slug); l++ {
			for m := 0; m < len(slug); m++ {
				compSlug := slug[k] + "-vs-" + slug[l] + "-vs-" + slug[m]
				if db.Conn.Where(&models.Comparison{Slug: compSlug}).First(&comparison).Error == nil {
					return true, compSlug
				}
			}
		}
	}

	compSlug := slug[0] + "-vs" + slug[1] + "-vs-" + slug[2]
	return false, compSlug
}
