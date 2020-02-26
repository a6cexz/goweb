package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"sync"

	"github.com/pkg/errors"
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
	s.Posts = s.loadPosts()
	s.Templates = s.loadTemplates()
}

func (s *BlogServer) lock() {
	s.mu.Lock()
}

func (s *BlogServer) unLock() {
	s.mu.Unlock()
}

func (s *BlogServer) loadPosts() map[int]*BlogPost {
	r := make(map[int]*BlogPost)
	r[0] = &BlogPost{
		ID:      0,
		Title:   "Title1",
		Date:    "21 Feb 2020",
		Link:    "https://google/link1",
		Content: "Test content1",
	}
	r[1] = &BlogPost{
		ID:      0,
		Title:   "Title2",
		Date:    "22 Feb 2020",
		Link:    "https://google/link2",
		Content: "Test content2",
	}
	return r
}

func (s *BlogServer) getPostByID(postID string) (*BlogPost, error) {
	s.lock()
	defer s.unLock()

	id, err := strconv.Atoi(postID)
	if err != nil {
		return nil, err
	}

	if blogPost, ok := s.Posts[id]; ok {
		return blogPost, nil

	}
	return nil, errors.Errorf("No post found for: %v", id)
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

func (s *BlogServer) addPost(post *BlogPost) {
	s.lock()
	defer s.unLock()
	s.Posts[post.ID] = post
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

	err := t.ExecuteTemplate(wr, "page", s)
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
