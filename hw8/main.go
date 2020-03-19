package main

import (
	_ "hw8/routers"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	logToFile := beego.AppConfig.DefaultBool("logToFile", false)
	if logToFile {
		logFileName := beego.AppConfig.String("logFileName")
		if err := logs.SetLogger(logs.AdapterFile, `{"filename":"`+logFileName+`"}`); err != nil {
			log.Print(err)
		}
	}

	beego.Info("Starting blog server")
	beego.Run()
}
