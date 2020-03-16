package controllers

import (
	"database/sql"
	"fmt"
	"hw5_2/models"
	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/pkg/errors"
)

// MainController controller
type MainController struct {
	beego.Controller
	DB *sql.DB
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
		id, err := strconv.Atoi(postID)
		if err != nil {
			err := errors.Wrapf(err, "Can not parse id value: %v", postID)
			http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
			log.Print(err)
			return
		}

		p := &models.BlogPost{}
		p.ID = id
		p.Title = req.FormValue("title")
		p.Date = req.FormValue("date")
		p.Link = req.FormValue("link")
		p.Content = req.FormValue("content")
		err = c.addPost(p)
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

func (c *MainController) getAllPosts() ([]models.BlogPost, error) {
	r := make([]models.BlogPost, 0)
	rows, err := c.DB.Query("select * from blog.posts")
	if err != nil {
		return r, err
	}
	defer rows.Close()

	for rows.Next() {
		post := models.BlogPost{}
		err := rows.Scan(&post.ID, &post.Title, &post.Date, &post.Link, &post.Content)
		if err != nil {
			log.Println(err)
			continue
		}
		r = append(r, post)
	}
	return r, nil
}

func (c *MainController) addPost(post *models.BlogPost) error {
	cmd := fmt.Sprintf("insert into blog.posts (title, postdate, link, content) values(\"%v\", \"%v\", \"%v\", \"%v\")", post.Title, post.Date, post.Link, post.Content)
	res, err := c.DB.Exec(cmd)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int(id)
	return nil
}

func (c *MainController) getPostByID(postID string) (models.BlogPost, error) {
	post := models.BlogPost{}

	row := c.DB.QueryRow(fmt.Sprintf("select * from blog.posts where posts.id = %v", postID))
	err := row.Scan(&post.ID, &post.Title, &post.Date, &post.Link, &post.Content)
	if err != nil {
		return post, err
	}

	return post, nil
}
