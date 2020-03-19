package controllers

import (
	ctx "context"
	"hw8/models"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MainController controller
type MainController struct {
	beego.Controller
	DB     *mongo.Client
	DBName string
}

func parseObjectID(str string) (primitive.ObjectID, error) {
	s := strings.TrimPrefix(str, "ObjectID(\"")
	hex := strings.TrimSuffix(s, "\")")
	return primitive.ObjectIDFromHex(hex)
}

// ListPosts gets main page
func (c *MainController) ListPosts() {
	beego.Info("ListPosts")

	posts, err := c.GetAllPosts()
	if err != nil {
		err = errors.Wrap(err, "Can not load posts")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		beego.Error(err)
		return
	}

	beego.Info("Loaded %v posts", len(posts))

	c.Data["Title"] = "Blog"
	c.Data["Posts"] = posts
	c.TplName = "index.tpl"
}

// ReadPost shows post
func (c *MainController) ReadPost() {
	beego.Info("ReadPost")

	req := c.Ctx.Request

	postID := req.URL.Query().Get("id")
	post, err := c.GetPostByID(postID)
	if err != nil {
		err := errors.Wrap(err, "No post found")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		beego.Error(err)
		return
	}

	beego.Info("Post loaded: %v", post.Title)

	c.Data["Title"] = post.Title
	c.Data["Post"] = post
	c.TplName = "post.tpl"
}

// EditPost shows post
func (c *MainController) EditPost() {
	beego.Info("EditPost")

	req := c.Ctx.Request

	postID := req.URL.Query().Get("id")
	post, err := c.GetPostByID(postID)
	if err != nil {
		err := errors.Wrap(err, "No post found")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		beego.Error(err)
		return
	}

	beego.Info("Edit post: %v", post.Title)

	c.Data["Title"] = post.Title
	c.Data["Post"] = post
	c.TplName = "editPost.tpl"
}

// UpdatePost updates data
func (c *MainController) UpdatePost() {
	beego.Info("UpdatePost")

	req := c.Ctx.Request
	postID := req.FormValue("id")
	if len(postID) > 0 {
		post := &models.BlogPost{}
		objID, err := parseObjectID(postID)
		if err != nil {
			err = errors.Wrap(err, "Can not parse post id")
			http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
			beego.Error(err)
			return
		}

		post.ID = objID
		post.Title = req.FormValue("title")
		post.Date = req.FormValue("date")
		post.Link = req.FormValue("link")
		post.Content = req.FormValue("content")
		err = c.UpdateBlogPost(post)
		if err != nil {
			err = errors.Wrap(err, "Can not create post")
			http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
			beego.Error(err)
			return
		}

		beego.Info("Updated post: %v", post.Title)

		c.Data["Title"] = post.Title
		c.Data["Post"] = post
		c.TplName = "post.tpl"
	}
}

// NewPost creates new post
func (c *MainController) NewPost() {
	beego.Info("NewPost")

	post, err := c.CreateNewPost(c.Ctx.ResponseWriter)
	if err != nil {
		err = errors.Wrap(err, "Can not create new post")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		beego.Error(err)
		return
	}

	if post != nil {
		beego.Info("New Post Created")
		c.Data["Title"] = post.Title
		c.Data["Post"] = post
		c.TplName = "post.tpl"
	}
}

// GetAllPosts gets all posts
func (c *MainController) GetAllPosts() ([]models.BlogPost, error) {
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

// GetPostByID gets post by id
func (c *MainController) GetPostByID(postID string) (*models.BlogPost, error) {
	col := c.DB.Database(c.DBName).Collection("posts")

	objID, err := parseObjectID(postID)
	if err != nil {
		err = errors.Wrap(err, "Can not parse id")
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		beego.Error(err)
		return nil, err
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	res := col.FindOne(ctx.TODO(), filter)
	post := &models.BlogPost{}
	err = res.Decode(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// AddPost new post
func (c *MainController) AddPost(post *models.BlogPost) error {
	col := c.DB.Database(c.DBName).Collection("posts")
	result, err := col.InsertOne(ctx.TODO(), post)
	if err != nil {
		return err
	}
	objID := result.InsertedID.(primitive.ObjectID)
	post.ID = objID
	return nil
}

// UpdateBlogPost updates post
func (c *MainController) UpdateBlogPost(post *models.BlogPost) error {
	col := c.DB.Database(c.DBName).Collection("posts")

	filter := bson.M{"_id": bson.M{"$eq": post.ID}}
	update := bson.M{"$set": bson.M{"title": post.Title, "link": post.Link, "date": post.Date, "content": post.Content}}

	_, err := col.UpdateOne(ctx.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// CreateTestPost creates new test post
func (c *MainController) CreateTestPost(wr http.ResponseWriter) {
	c.CreateNewPost(wr)
}

// CreateNewPost creates new post
func (c *MainController) CreateNewPost(wr http.ResponseWriter) (*models.BlogPost, error) {
	post := &models.BlogPost{}
	post.Title = "TestPost1"
	post.Date = "2019-10-01"
	post.Link = "TestLink"
	post.Content = "TestContent"
	err := c.AddPost(post)
	if err != nil {
		err = errors.Wrap(err, "Can not create post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		beego.Error(err)
		return nil, err
	}
	return post, nil
}
