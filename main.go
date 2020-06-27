package main

import (
	"Parser-Golang/db"
	_ "Parser-Golang/routers"
	"Parser-Golang/services"
	"Parser-Golang/utilities"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/carlescere/scheduler"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
)

func main() {
	RegisterTemplateFuncs()
	conn := getDBConnection()
	defer conn.Close()
	if getProxyConnection() {
		RegisterSchedulerFuncs()
	}

	beego.Run("localhost:8080")
}

func RegisterSchedulerFuncs() {
	scheduler.Every(5).Hours().Run(services.AddNewComparisons)
	//scheduler.Every().Sunday().At("02:30").Run(services.AddNewComparisons)
	//scheduler.Every(1).Hours().Run(services.DeleteComparisonIfHasProblem)
	//scheduler.Every().Day().At("00:01").Run(services.UpdateGitData)
}

func getProxyConnection() bool {
	var user = "DLcGHY"
	var password = "ucAWpv"
	var address = "91.188.242.138"
	var port = "9180"
	var proxy = fmt.Sprintf("http://%s:%s@%s:%s", user, password, address, port)

	os.Setenv("HTTP_PROXY", proxy)
	_, err := http.Get("https://stackshare.io")

	if err != nil {
		logs.Warn("Proxy Error!")
		return false
	}

	return true
}

func getDBConnection() *gorm.DB {
	dbUser := "osmanlie_hakan"
	dbPass := "13399771Hm"
	dbName := "osmanlie_local"
	dbHost := "192.185.73.49"

	conn, err := db.Connection(dbUser, dbPass, dbName, dbHost)
	if err != nil {
		panic(err.Error())
	}

	return conn
}

func RegisterTemplateFuncs() {
	beego.AddFuncMap("IfNotZero", utilities.IfNotZero)
	beego.AddFuncMap("GetFirstWord", utilities.GetFirstWord)
	beego.AddFuncMap("PreloadStackPros", utilities.PreloadStackPros)
	beego.AddFuncMap("PreloadStackCons", utilities.PreloadStackCons)
	beego.AddFuncMap("PreloadStackCompanies", utilities.PreloadStackCompanies)
}
