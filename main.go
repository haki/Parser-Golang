package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"parser/db"
	_ "parser/routers"
	"parser/services"
)

func main() {
	getProxyConnection()
	conn := getDBConnection()
	defer conn.Close()

	//scheduler.Every().Sunday().At("08:30").Run(UpdateComparison)

	beego.Run("localhost:8080")
}

func UpdateComparison() {
	services.UpdateData()
}

func getProxyConnection() {
	var user = "DLcGHY"
	var password = "ucAWpv"
	var address = "91.188.242.138"
	var port = "9180"
	var proxy = fmt.Sprintf("http://%s:%s@%s:%s", user, password, address, port)

	os.Setenv("HTTP_PROXY", proxy)
	_, err := http.Get("https://stackshare.io")

	if err != nil {
		logs.Error("Proxy Error")
	}
}

func getDBConnection() *gorm.DB {
	dbUser := "user"
	dbPass := "password"
	dbName := "db"
	dbHost := "localhost"

	conn, err := db.Connection(dbUser, dbPass, dbName, dbHost)
	if err != nil {
		panic(err.Error())
	}

	return conn
}
