package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"codeberg.org/mislavzanic/main/internal/posts"
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
	p := posts.ReadBlogPost(fmt.Sprintf("%s/%s.md", POSTSDIR, pageId))

	posts.RenderFromTemplate(w, "post.html", fmt.Sprintf("%s/post.html", HTMLDIR), template.FuncMap{"toURL": posts.GetUrl("blog"), "markDown": posts.ToMarkdown, "afterEpoch": posts.AfterEpoch}, p)
}

func AboutSection(w http.ResponseWriter, req *http.Request) {
	p := posts.ReadBlogPost(fmt.Sprintf("%s/about.md", ABOUTDIR))
	posts.RenderFromTemplate(w, "post.html", fmt.Sprintf("%s/post.html", HTMLDIR), template.FuncMap{"markDown": posts.ToMarkdown, "afterEpoch": posts.AfterEpoch}, p)
}

func ViewProjects(w http.ResponseWriter, req *http.Request) {
	allPosts := posts.GetAllPosts(PROJDIR)
	posts.RenderFromTemplate(w, "index.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": posts.GetUrl("projects"), "markDown": posts.ToMarkdown, "afterEpoch": posts.AfterEpoch}, allPosts)
}

func GetProject(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["projId"]
	p := posts.ReadBlogPost(fmt.Sprintf("%s/%s.md", PROJDIR, id))

	posts.RenderFromTemplate(w, "post.html", fmt.Sprintf("%s/post.html", HTMLDIR), template.FuncMap{"markDown": posts.ToMarkdown, "afterEpoch": posts.AfterEpoch}, p)
}

func FilterByTag(w http.ResponseWriter, req *http.Request) {
	tagId := mux.Vars(req)["tagId"]
	allPosts := posts.FindBlogPosts(tagId, POSTSDIR)

	posts.RenderFromTemplate(w, "tags.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": posts.GetUrl("blog"), "markDown": posts.ToMarkdown, "afterEpoch": posts.AfterEpoch}, allPosts)
}

func ViewAllPosts(w http.ResponseWriter, req *http.Request) {
	allPosts := posts.GetAllPosts(POSTSDIR)
	posts.RenderFromTemplate(w, "index.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": posts.GetUrl("blog"), "markDown": posts.ToMarkdown, "afterEpoch": posts.AfterEpoch}, allPosts)
}

