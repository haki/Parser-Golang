package models

import "github.com/jinzhu/gorm"

type Stack struct {
	gorm.Model

	Name        string       `gorm:"type:varchar(100)"`
	Slug        string       `gorm:"type:varchar(100);not null;unique"`
	Description string       `gorm:"type:text;"`
	Image       string       `gorm:"type:varchar(255)"`
	Website     string       `gorm:"type:varchar(255)"`
	View        int64        `gorm:"type:int;default:0;"`
	GitUrl      string       `gorm:"type:varchar(255)"`
	Fork        int          `gorm:"type:varchar(150)"`
	Star        int          `gorm:"type:varchar(150)"`
	Watch       int          `gorm:"type:varchar(150)"`
	Comparisons []Comparison `gorm:"many2many:comparison_stacks"`
	Companies   []Company    `gorm:"many2many:stack_companies"`
	Cons        []Cons       `gorm:"many2many:stack_cons"`
	Pros        []Pros       `gorm:"many2many:stack_pros"`
}

func (Stack) TableName() string {
	return "stacks"
}
