package models

import "github.com/jinzhu/gorm"

type Pros struct {
	gorm.Model

	Text    string `gorm:"type:varchar(250)"`
	Point   int    `gorm:"type:int"`
	Enabled bool
	Stacks  Stack `gorm:"many2many:stack_pros"`
}

func (Pros) TableName() string {
	return "pros"
}
