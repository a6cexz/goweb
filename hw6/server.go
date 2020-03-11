package main

import (
	ctx "context"
	"html/template"
	"log"
	"net/http"
	"path"
	"sort"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	editPost  = "editPost.html"
	showPost  = "showPost.html"
	listPosts = "listPosts.html"
)

var templateFiles = []string{
	editPost,
	showPost,
	listPosts,
}

// BlogServer struct
type BlogServer struct {
	mu        sync.Mutex
	client    *mongo.Client
	dbName    string
	Title     string
	Posts     map[string]BlogPost
	Templates map[string]*template.Template
}

// BlogPost struct
type BlogPost struct {
	ID      string
	Title   string
	Date    string
	Link    string
	Content string
}

func main() {
	server := BlogServer{}
	server.init()
	defer server.shutDown()

	beego.Get("/", server.handleRoot)
	beego.Post("/", server.handleEditPost)
	beego.Get("/post/", server.handlePost)
	beego.Get("/edit/", server.handleEdit)
	beego.Get("/new/", server.handleNew)

	beego.Run()
}

func (s *BlogServer) init() {
	s.Title = "Blog"
	s.Templates = s.loadTemplates()
	s.dbName = "BlogData"

	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	s.client = db

	err = s.client.Connect(ctx.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func (s *BlogServer) loadTemplates() map[string]*template.Template {
	r := make(map[string]*template.Template, len(templateFiles))
	for _, name := range templateFiles {
		t := template.Must(template.New("MyTemplate").ParseFiles(path.Join("templates", string(name))))
		r[name] = t
	}
	return r
}

func (s *BlogServer) shutDown() {
	s.client.Disconnect(ctx.TODO())
}

func (s *BlogServer) handleRoot(ctx *context.Context) {
	t := s.getTemplate(listPosts)
	if t == nil {
		err := errors.Errorf("No template found: %v", listPosts)
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	posts, err := s.getAllPosts()
	if err != nil {
		err = errors.Wrap(err, "Can not load posts")
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if len(posts) == 0 {
		s.createTestPost(ctx.ResponseWriter)
		ctx.Redirect(http.StatusOK, "/")
		return
	}

	s.Posts = posts
	err = t.ExecuteTemplate(ctx.ResponseWriter, "page", s)
	if err != nil {
		err = errors.Wrap(err, "Can not execute template")
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *BlogServer) createTestPost(wr http.ResponseWriter) {
	s.createNewPost(wr, "0")
}

func (s *BlogServer) createNewPost(wr http.ResponseWriter, postID string) {
	p1 := &BlogPost{}
	p1.ID = postID
	p1.Title = "TestPost1"
	p1.Date = "2019-10-01"
	p1.Link = "TestLink"
	p1.Content = "TestContent"
	err := s.addPost(p1)
	if err != nil {
		err = errors.Wrap(err, "Can not create test post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
	}
}

func (s *BlogServer) handlePost(ctx *context.Context) {
	t := s.getTemplate(showPost)
	if t == nil {
		err := errors.Errorf("No template found: %v", showPost)
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	postID := ctx.Request.URL.Query().Get("id")
	post, err := s.getPostByID(postID)
	if err != nil {
		err := errors.Wrap(err, "No post found")
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	err = t.ExecuteTemplate(ctx.ResponseWriter, "page", post)
	if err != nil {
		err = errors.Wrap(err, "Can not execute template")
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *BlogServer) handleEdit(ctx *context.Context) {
	t := s.getTemplate(editPost)
	if t == nil {
		err := errors.Errorf("No template found: %v", editPost)
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	postID := ctx.Request.URL.Query().Get("id")
	post, err := s.getPostByID(postID)
	if err != nil {
		err := errors.Wrap(err, "No post found")
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	err = t.ExecuteTemplate(ctx.ResponseWriter, "page", post)
	if err != nil {
		err = errors.Wrap(err, "Can not execute template")
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
	}
}

func (s *BlogServer) handleEditPost(ctx *context.Context) {
	s.editExistingPost(ctx.ResponseWriter, ctx.Request)
	ctx.Redirect(http.StatusOK, "/")
}

func (s *BlogServer) handleNew(ctx *context.Context) {
	posts, err := s.getAllPosts()
	if err != nil {
		err = errors.Wrap(err, "Can not load posts")
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if len(posts) == 0 {
		s.createTestPost(ctx.ResponseWriter)
		ctx.Redirect(http.StatusOK, "/")
		return
	}

	ids := []string{}
	for _, p := range posts {
		ids = append(ids, p.ID)
	}
	sort.Strings(ids)
	last := ids[len(ids)-1]
	mewID := last + "1"
	s.createNewPost(ctx.ResponseWriter, mewID)
}

func (s *BlogServer) getTemplate(name string) *template.Template {
	s.mu.Lock()
	defer s.mu.Unlock()

	if template, ok := s.Templates[name]; ok {
		return template
	}
	return nil
}

func (s *BlogServer) editExistingPost(wr http.ResponseWriter, req *http.Request) {
	postID := req.FormValue("id")
	if len(postID) > 0 {
		p := &BlogPost{}
		p.ID = postID
		p.Title = req.FormValue("title")
		p.Date = req.FormValue("date")
		p.Link = req.FormValue("link")
		p.Content = req.FormValue("content")
		err := s.updatePost(p)
		if err != nil {
			err = errors.Wrap(err, "Can not update post")
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s *BlogServer) getAllPosts() (map[string]BlogPost, error) {
	r := make(map[string]BlogPost)

	c := s.client.Database(s.dbName).Collection("posts")

	cur, err := c.Find(ctx.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	posts := []BlogPost{}
	err = cur.All(ctx.TODO(), &posts)

	if err != nil {
		return nil, err
	}

	for _, p := range posts {
		r[p.ID] = p
	}
	return r, nil
}

func (s *BlogServer) getPostByID(postID string) (*BlogPost, error) {
	c := s.client.Database(s.dbName).Collection("posts")
	filter := bson.M{"id": bson.M{"$eq": postID}}
	res := c.FindOne(ctx.TODO(), filter)
	post := &BlogPost{}
	err := res.Decode(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *BlogServer) addPost(post *BlogPost) error {
	c := s.client.Database(s.dbName).Collection("posts")
	_, err := c.InsertOne(ctx.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

func (s *BlogServer) updatePost(post *BlogPost) error {
	c := s.client.Database(s.dbName).Collection("posts")

	filter := bson.M{"id": bson.M{"$eq": post.ID}}
	update := bson.M{"$set": bson.M{"title": post.Title, "link": post.Link, "date": post.Date, "content": post.Content}}

	_, err := c.UpdateOne(ctx.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
