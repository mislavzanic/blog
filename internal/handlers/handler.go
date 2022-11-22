package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"codeberg.org/mislavzanic/main/internal/posts"
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



func (s Site) PageHandler(w http.ResponseWriter, req *http.Request) {
	pageId := mux.Vars(req)["pageId"]
	p, err := s.findPost(fmt.Sprintf("%s/%s.md", POSTSDIR, pageId))

	if err != nil {
		log.Fatal(err)
	}

	posts.RenderFromTemplate(
		w,
		"post.html",
		fmt.Sprintf("%s/post.html", HTMLDIR),
		template.FuncMap{
			"toURL": posts.GetUrl("blog"),
			"markDown": posts.ToMarkdown,
			"afterEpoch": posts.AfterEpoch,
		},
		p,
	)
}

func (s Site) AboutSection(w http.ResponseWriter, req *http.Request) {
	posts.RenderFromTemplate(w, "post.html", fmt.Sprintf("%s/post.html", HTMLDIR), template.FuncMap{"markDown": posts.ToMarkdown, "afterEpoch": posts.AfterEpoch}, s.About)
}

func (s Site) ViewProjects(w http.ResponseWriter, req *http.Request) {
	posts.RenderFromTemplate(
		w,
		"index.html",
		fmt.Sprintf("%s/index.html", HTMLDIR),
		template.FuncMap{
			"toURL": posts.GetUrl("projects"),
			"markDown": posts.ToMarkdown,
			"afterEpoch": posts.AfterEpoch,
		},
		s.Projects,
	)
}

func (s Site) FilterByTag(w http.ResponseWriter, req *http.Request) {
	tagId := mux.Vars(req)["tagId"]
	allPosts := posts.FindBlogPosts(tagId, POSTSDIR)

	posts.RenderFromTemplate(w, "tags.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": posts.GetUrl("blog"), "markDown": posts.ToMarkdown, "afterEpoch": posts.AfterEpoch}, allPosts)
}

func (s Site) ViewAllPosts(w http.ResponseWriter, req *http.Request) {
	posts.RenderFromTemplate(
		w, "index.html",
		fmt.Sprintf("%s/index.html", HTMLDIR),
		template.FuncMap{
			"toURL": posts.GetUrl("blog"),
			"markDown": posts.ToMarkdown,
			"afterEpoch": posts.AfterEpoch,
		},
		s.Blog,
	)
}

