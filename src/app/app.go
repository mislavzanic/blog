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

type Site struct {
	Blog     posts.Posts
	About    posts.Page
	Projects posts.Posts
}

const (
	BLOGDIR   = "blog/posts"
	ABOUTPAGE = "blog/about/about.md"
	PROJDIR   = "blog/projects"
)

func LoadSite() Site {
	return Site{
		Blog: posts.GetAllPosts(BLOGDIR),
		About: *posts.ReadBlogPost(ABOUTPAGE),
		Projects: posts.GetAllPosts(PROJDIR),
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
	renderer.RenderFromTemplate(
		w,
		"index.html",
		"html/index.html",
		template.FuncMap{
			"toURL": renderer.GetUrl(index.Uri),
			"markDown": renderer.ToMarkdown,
			"afterEpoch": renderer.AfterEpoch,
		},
		index,
	)
}

func (s Site) RenderFilterIndex(w http.ResponseWriter, index posts.Posts, filter string) {
	s.RenderIndex(w, index.FindBlogPosts(filter))
}

func (s Site) RenderPage(w http.ResponseWriter, path string) {
	post, err := s.findPost(path)

	if err != nil {
		log.Fatal(err)
	}

	renderer.RenderFromTemplate(
		w,
		"post.html",
		"html/post.html",
		template.FuncMap{
			"toURL": renderer.GetUrl("posts"),
			"markDown": renderer.ToMarkdown,
			"afterEpoch": renderer.AfterEpoch,
		},
		post,
	)
}
