package models

import "github.com/jinzhu/gorm"

type Company struct {
	gorm.Model

	Name        string  `gorm:"type:varchar(100);not null"`
	Slug        string  `gorm:"type:varchar(100);not null;unique"`
	Image       string  `gorm:"type:varchar(255);"`
	Website     string  `gorm:"type:varchar(255);"`
	Email       string  `gorm:"type:varchar(300);"`
	Github      string  `gorm:"type:varchar(255)"`
	LinkedIn    string  `gorm:"type:varchar(255)"`
	Facebook    string  `gorm:"type:varchar(255)"`
	Description string  `gorm:"type:varchar(500)"`
	Country     string  `gorm:"type:varchar(50)"`
	Stacks      []Stack `gorm:"many2many:stack_companies;"`
}

func (Company) TableName() string {
	return "companies"
}
