package models

import "github.com/jinzhu/gorm"

type Cons struct {
	gorm.Model

	Text    string `gorm:"type:varchar(250)"`
	Point   int    `gorm:"type:int"`
	Enabled bool
	Stacks  Stack `gorm:"many2many:stack_cons"`
}

func (Cons) TableName() string {
	return "cons"
}
