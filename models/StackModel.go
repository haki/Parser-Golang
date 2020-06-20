package models

import "github.com/jinzhu/gorm"

type Stack struct {
	gorm.Model

	Name        string       `gorm:"type:varchar(100)"`
	Slug        string       `gorm:"type:varchar(100);not null;unique"`
	Description string       `gorm:"type:varchar(1000);"`
	Image       string       `gorm:"type:varchar(450)"`
	GithubUrl   string       `gorm:"type:varchar(350)"`
	Fork        string       `gorm:"type:varchar(150)"`
	Star        string       `gorm:"type:varchar(150)"`
	Watch       string       `gorm:"type:varchar(150)"`
	Comparisons []Comparison `gorm:"many2many:comparison_stacks"`
	Companies   []Company    `gorm:"many2many:stack_companies"`
	Cons        []Cons       `gorm:"many2many:stack_cons"`
	Pros        []Pros       `gorm:"many2many:stack_pros"`
}

func (Stack) TableName() string {
	return "stacks"
}
