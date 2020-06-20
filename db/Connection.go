package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"parser/models"
)

var Conn *gorm.DB

func Connection(dbUser string, dbPass string, dbName string, dbHost string) (db *gorm.DB, err error) {
	dbConnString := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err = gorm.Open("mysql", dbConnString)
	Conn = db

	if err == nil {
		Conn.AutoMigrate(&models.Company{}, &models.Comparison{}, &models.Cons{}, &models.Pros{}, &models.Stack{})
	}

	return
}
