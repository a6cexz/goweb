package routers

import (
	ctx "context"
	"hw8/controllers"
	"log"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	beego.Info("Starting db")
	dbURI := beego.AppConfig.String("dbUri")
	db, err := mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		beego.Critical(err)
		log.Fatal(err)
	}

	err = db.Connect(ctx.TODO())
	if err != nil {
		beego.Critical(err)
		log.Fatal(err)
	}

	dbName := beego.AppConfig.String("dbName")
	controller := &controllers.MainController{DB: db, DBName: dbName}
	beego.Router("/", controller, "get:ListPosts")
	beego.Router("/post", controller, "get:ReadPost")
	beego.Router("/edit", controller, "get:EditPost")
	beego.Router("/edit", controller, "post:UpdatePost")
	beego.Router("/new", controller, "post:NewPost")
}
