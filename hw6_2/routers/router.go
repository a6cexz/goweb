package routers

import (
	ctx "context"
	"hw6/controllers"
	"log"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Connect(ctx.TODO())
	if err != nil {
		log.Fatal(err)
	}

	controller := &controllers.MainController{DB: db, DBName: "BlogData"}
	beego.Router("/", controller)
	beego.Router("/post", controller, "get:ShowPost")
	beego.Router("/edit", controller, "get:EditPost")
	beego.Router("/new", controller, "get:NewPost")
}
