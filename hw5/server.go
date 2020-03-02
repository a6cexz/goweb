package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"sync"

	"database/sql"

	_ "github.com/go-sql-driver/MySQL"
	"github.com/pkg/errors"
)

var database *sql.DB

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
	Title     string
	Posts     map[int]*BlogPost
	Templates map[string]*template.Template
}

// BlogPost struct
type BlogPost struct {
	ID      int
	Title   string
	Date    string
	Link    string
	Content string
}

func main() {
	db, err := sql.Open("MySQL", "root:pwd1234567@/blog")
	if err != nil {
		log.Fatal(err)
	}
	database = db
	defer database.Close()

	router := http.NewServeMux()

	server := BlogServer{}
	server.init()

	router.HandleFunc("/", server.handleRoot)
	router.HandleFunc("/post/", server.handlePost)
	router.HandleFunc("/edit/", server.handleEdit)

	port := "8080"
	log.Printf("start server on port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func (s *BlogServer) init() {
	s.Title = "Blog"
	s.Templates = s.loadTemplates()
}

func (s *BlogServer) lock() {
	s.mu.Lock()
}

func (s *BlogServer) unLock() {
	s.mu.Unlock()
}

func (s *BlogServer) getAllPosts() (map[int]*BlogPost, error) {
	r := make(map[int]*BlogPost)
	rows, err := database.Query("select * from blog.posts")
	if err != nil {
		return r, err
	}
	defer rows.Close()

	for rows.Next() {
		post := &BlogPost{}
		err := rows.Scan(&post.ID, &post.Title, &post.Date, &post.Link, &post.Content)
		if err != nil {
			log.Println(err)
			continue
		}
		r[post.ID] = post
	}
	return r, nil
}

func (s *BlogServer) getPostByID(postID string) (*BlogPost, error) {
	post := &BlogPost{}

	row := database.QueryRow(fmt.Sprintf("select * from blog.posts where posts.id = %v", postID))
	err := row.Scan(&post.ID, &post.Title, &post.Date, &post.Link, &post.Content)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (s *BlogServer) getTemplate(name string) *template.Template {
	s.lock()
	defer s.unLock()

	if template, ok := s.Templates[name]; ok {
		return template
	}
	return nil
}

func (s *BlogServer) loadTemplates() map[string]*template.Template {
	r := make(map[string]*template.Template, len(templateFiles))
	for _, name := range templateFiles {
		t := template.Must(template.New("MyTemplate").ParseFiles(path.Join("templates", string(name))))
		r[name] = t
	}
	return r
}

func (s *BlogServer) createNewPost(wr http.ResponseWriter, req *http.Request) {
	postID := req.FormValue("id")
	if len(postID) > 0 {
		id, err := strconv.Atoi(postID)
		if err != nil {
			err := errors.Wrapf(err, "Can not parse id value: %v", postID)
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			log.Print(err)
			return
		}

		p := &BlogPost{}
		p.ID = id
		p.Title = req.FormValue("title")
		p.Date = req.FormValue("date")
		p.Link = req.FormValue("link")
		p.Content = req.FormValue("content")
		s.addPost(p)
	}
}

func (s *BlogServer) addPost(post *BlogPost) error {
	cmd := fmt.Sprintf("insert into blog.posts (title, postdate, link, content) values(\"%v\", \"%v\", \"%v\", \"%v\")", post.Title, post.Date, post.Link, post.Content)
	res, err := database.Exec(cmd)
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

func (s *BlogServer) handleRoot(wr http.ResponseWriter, req *http.Request) {
	t := s.getTemplate(listPosts)
	if t == nil {
		err := errors.Errorf("No template found: %v", listPosts)
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	s.createNewPost(wr, req)

	posts, err := s.getAllPosts()
	if err != nil {
		err = errors.Wrap(err, "Can not load posts")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	s.Posts = posts
	err = t.ExecuteTemplate(wr, "page", s)
	if err != nil {
		err = errors.Wrap(err, "Can not execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *BlogServer) handlePost(wr http.ResponseWriter, req *http.Request) {
	t := s.getTemplate(showPost)
	if t == nil {
		err := errors.Errorf("No template found: %v", showPost)
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	postID := req.URL.Query().Get("id")
	post, err := s.getPostByID(postID)
	if err != nil {
		err := errors.Wrap(err, "No post found")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	err = t.ExecuteTemplate(wr, "page", post)
	if err != nil {
		err = errors.Wrap(err, "Can not execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *BlogServer) handleEdit(wr http.ResponseWriter, req *http.Request) {
	t := s.getTemplate(editPost)
	if t == nil {
		err := errors.Errorf("No template found: %v", editPost)
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	postID := req.URL.Query().Get("id")
	post, err := s.getPostByID(postID)
	if err != nil {
		err := errors.Wrap(err, "No post found")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	err = t.ExecuteTemplate(wr, "page", post)
	if err != nil {
		err = errors.Wrap(err, "Can not execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
