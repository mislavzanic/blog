package app

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mislavzanic/blog/src/app/posts"
	"github.com/mislavzanic/blog/src/app/renderer"
)

type void struct{}
type set map[string]void

type Site struct {
	Blog     posts.Posts
	About    posts.Page
	Projects posts.Posts
}

const (
	HTMLDIR  = "static/html"
	CSSDIR   = "static/css"
	JSDIR    = "static/js"
	BLOGDIR   = "webContent/blog"
	ABOUTPAGE = "webContent/about/about.md"
	PROJDIR   = "webContent/projects"
)

func LoadSite() Site {
	return Site{
		Blog: posts.GetAllPosts("blog"),
		About: *posts.ReadBlogPost(ABOUTPAGE),
		Projects: posts.GetAllPosts("projects"),
	}
}

func (s Site) findPost(path string) (*posts.Page, error) {

	if path == ABOUTPAGE {
		return &s.About, nil
	}

	for _, post := range s.Blog.Pages {
		if post.Path == path {
			return post, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Can't find post at %s", path))
}

func (s Site) RenderIndex(w http.ResponseWriter, index posts.Posts) {
	render(w, "index.html", index.Uri, index)
}

func (s Site) RenderFilterIndex(w http.ResponseWriter, index posts.Posts, filter string) {
	s.RenderIndex(w, index.FindBlogPosts(filter))
}

func (s Site) RenderPage(w http.ResponseWriter, path string) {
	post, err := s.findPost(path)

	if err != nil {
		log.Fatal(err)
	}

	render(w, "post.html", "blog", post)
}

func render(w http.ResponseWriter, templateName, uri string, data interface{}) {
	renderer.RenderFromTemplate(
		w,
		templateName,
		[]string{
			fmt.Sprintf("%s/%s", HTMLDIR, templateName),
			fmt.Sprintf("%s/header.html", HTMLDIR),
			fmt.Sprintf("%s/footer.html", HTMLDIR),
			fmt.Sprintf("%s/head.html", HTMLDIR),
		},
		template.FuncMap{
			"toURL": renderer.GetUrl(uri),
			"markDown": renderer.ToMarkdown,
			"afterEpoch": renderer.AfterEpoch,
		},
		data,
	)
}
