package routers

import (
	"database/sql"
	"hw5_2/controllers"
	"log"

	"github.com/astaxie/beego"
	// sql driver
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db, err := sql.Open("mysql", "root:root@/blog")
	if err != nil {
		log.Fatal(err)
		return
	}

	controller := &controllers.MainController{DB: db}
	beego.Router("/", controller)
	beego.Router("/post", controller, "get:ShowPost")
	beego.Router("/edit", controller, "get:EditPost")
}
