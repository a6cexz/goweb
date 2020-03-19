package tests

import (
	ctx "context"
	"hw8/controllers"
	_ "hw8/routers"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func startDb(dbName string) (*mongo.Client, error) {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.Connect(ctx.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db.Database(dbName).Collection("posts").Drop(ctx.TODO())
	return db, nil
}

func TestPosts(t *testing.T) {
	db, err := startDb("TestDb")
	if err != nil {
		t.Error(err)
		return
	}

	controller := &controllers.MainController{DB: db, DBName: "TestDb"}
	posts, err := controller.GetAllPosts()
	if err != nil {
		t.Error(err)
		return
	}

	if len(posts) != 0 {
		t.Error("Should be 0 posts")
		return
	}

	w := httptest.NewRecorder()
	newPost, err := controller.CreateNewPost(w)
	if w.Code != 200 {
		t.Error("Should 200")
		return
	}

	posts, err = controller.GetAllPosts()
	if err != nil {
		t.Error(err)
		return
	}
	if len(posts) != 1 {
		t.Error("Should be 1 posts")
		return
	}

	postIDHex := newPost.ID.Hex()
	post, err := controller.GetPostByID(postIDHex)
	if err != nil {
		t.Error(err)
		return
	}
	if post.Title != "TestPost1" {
		t.Error("Should be TestPost1 title")
		return
	}
}

// TestBeego is a sample to run an endpoint test
func TestGetRoot(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("Test Get /\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

func TestShowPost(t *testing.T) {
	r, _ := http.NewRequest("POST", "/new", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("Test new post", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}
