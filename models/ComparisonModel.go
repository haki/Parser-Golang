package models

import (
	"github.com/jinzhu/gorm"
)

type Comparison struct {
	gorm.Model

	Name       string  `gorm:"type:varchar(100);not null;"`
	Slug       string  `gorm:"type:varchar(100);not null;unique;"`
	View       int64   `gorm:"type:int;default:0;"`
	SourcePage string  `gorm:"type:varchar(200);"`
	Stacks     []Stack `gorm:"many2many:comparison_stacks;"`
}

func (Comparison) TableName() string {
	return "comparisons"
}
