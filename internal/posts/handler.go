package posts

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	HTMLDIR  = "html"
	CSSDIR   = "css"
	JSDIR    = "js"
	GITDIR   = "BlogPosts"
	POSTSDIR = "BlogPosts/posts"
	ABOUTDIR = "BlogPosts/about"
	PROJDIR  = "BlogPosts/projects"
)


func PageHandler(w http.ResponseWriter, req *http.Request) {
	pageId := mux.Vars(req)["pageId"]
	p := readBlogPost(fmt.Sprintf("%s/%s.md", POSTSDIR, pageId))

	renderFromTemplate(w, "post.html", fmt.Sprintf("%s/post.html", HTMLDIR), template.FuncMap{"markDown": markDowner, "afterEpoch": AfterEpoch}, p)
}

func AboutSection(w http.ResponseWriter, req *http.Request) {
	p := readBlogPost(fmt.Sprintf("%s/about.md", ABOUTDIR))
	renderFromTemplate(w, "post.html", fmt.Sprintf("%s/post.html", HTMLDIR), template.FuncMap{"markDown": markDowner, "afterEpoch": AfterEpoch}, p)
}

func ViewProjects(w http.ResponseWriter, req *http.Request) {
	posts := getAllPosts(PROJDIR)
	renderFromTemplate(w, "index.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": getUrl("projects"), "markDown": markDowner, "afterEpoch": AfterEpoch}, posts)
}

func GetProject(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["projId"]
	p := readBlogPost(fmt.Sprintf("%s/%s.md", PROJDIR, id))

	renderFromTemplate(w, "post.html", fmt.Sprintf("%s/post.html", HTMLDIR), template.FuncMap{"markDown": markDowner, "afterEpoch": AfterEpoch}, p)
}

func FilterByTag(w http.ResponseWriter, req *http.Request) {
	tagId := mux.Vars(req)["tagId"]
	posts := findBlogPosts(tagId)

	renderFromTemplate(w, "tags.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": getUrl("blog"), "markDown": markDowner, "afterEpoch": AfterEpoch}, posts)
}

func ViewAllPosts(w http.ResponseWriter, req *http.Request) {
	posts := getAllPosts(POSTSDIR)
	renderFromTemplate(w, "index.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": getUrl("blog"), "markDown": markDowner, "afterEpoch": AfterEpoch}, posts)
}

