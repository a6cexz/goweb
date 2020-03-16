package controllers

import (
	ctx "context"
	"hw6/models"
	"log"
	"net/http"
	"sort"

	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MainController controller
type MainController struct {
	beego.Controller
	DB     *mongo.Client
	DBName string
}

// Get gets main page
func (c *MainController) Get() {
	posts, err := c.getAllPosts()
	if err != nil {
		err = errors.Wrap(err, "Can not load posts")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	c.Data["Title"] = "Blog"
	c.Data["Posts"] = posts
	c.TplName = "index.tpl"
}

// Post updates data
func (c *MainController) Post() {
	req := c.Ctx.Request
	postID := req.FormValue("id")
	if len(postID) > 0 {
		p := &models.BlogPost{}
		p.ID = postID
		p.Title = req.FormValue("title")
		p.Date = req.FormValue("date")
		p.Link = req.FormValue("link")
		p.Content = req.FormValue("content")
		err := c.addPost(p)
		if err != nil {
			err = errors.Wrap(err, "Can not create post")
			http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
			log.Print(err)
			return
		}

		c.Data["Title"] = p.Title
		c.Data["Post"] = p
		c.TplName = "post.tpl"
	}
}

// ShowPost shows post
func (c *MainController) ShowPost() {
	req := c.Ctx.Request

	postID := req.URL.Query().Get("id")
	post, err := c.getPostByID(postID)
	if err != nil {
		err := errors.Wrap(err, "No post found")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	c.Data["Title"] = post.Title
	c.Data["Post"] = post
	c.TplName = "post.tpl"
}

// EditPost shows post
func (c *MainController) EditPost() {
	req := c.Ctx.Request

	postID := req.URL.Query().Get("id")
	post, err := c.getPostByID(postID)
	if err != nil {
		err := errors.Wrap(err, "No post found")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	c.Data["Title"] = post.Title
	c.Data["Post"] = post
	c.TplName = "editPost.tpl"
}

// NewPost creates new post
func (c *MainController) NewPost() {
	posts, err := c.getAllPosts()
	if err != nil {
		err = errors.Wrap(err, "Can not load posts")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if len(posts) == 0 {
		c.createTestPost(c.Ctx.ResponseWriter)
		c.Ctx.Redirect(http.StatusOK, "/")
		return
	}

	ids := []string{}
	for _, p := range posts {
		ids = append(ids, p.ID)
	}
	sort.Strings(ids)
	last := ids[len(ids)-1]
	mewID := last + "1"
	c.createNewPost(c.Ctx.ResponseWriter, mewID)
}

func (c *MainController) getAllPosts() ([]models.BlogPost, error) {
	col := c.DB.Database(c.DBName).Collection("posts")

	cur, err := col.Find(ctx.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	posts := []models.BlogPost{}
	err = cur.All(ctx.TODO(), &posts)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (c *MainController) getPostByID(postID string) (*models.BlogPost, error) {
	col := c.DB.Database(c.DBName).Collection("posts")
	filter := bson.M{"id": bson.M{"$eq": postID}}
	res := col.FindOne(ctx.TODO(), filter)
	post := &models.BlogPost{}
	err := res.Decode(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (c *MainController) addPost(post *models.BlogPost) error {
	col := c.DB.Database(c.DBName).Collection("posts")
	_, err := col.InsertOne(ctx.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

func (c *MainController) updatePost(post *models.BlogPost) error {
	col := c.DB.Database(c.DBName).Collection("posts")

	filter := bson.M{"id": bson.M{"$eq": post.ID}}
	update := bson.M{"$set": bson.M{"title": post.Title, "link": post.Link, "date": post.Date, "content": post.Content}}

	_, err := col.UpdateOne(ctx.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *MainController) createTestPost(wr http.ResponseWriter) {
	c.createNewPost(wr, "0")
}

func (c *MainController) createNewPost(wr http.ResponseWriter, postID string) {
	p1 := &models.BlogPost{}
	p1.ID = postID
	p1.Title = "TestPost1"
	p1.Date = "2019-10-01"
	p1.Link = "TestLink"
	p1.Content = "TestContent"
	err := c.addPost(p1)
	if err != nil {
		err = errors.Wrap(err, "Can not create test post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
	}
}
