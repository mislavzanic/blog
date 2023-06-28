package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mislavzanic/blog/src/app"
	"github.com/mislavzanic/blog/src/metrics"
)

const (
	HTMLDIR  = "html"
	CSSDIR   = "css"
	JSDIR    = "js"
)


func ViewProjects(site app.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		site.RenderIndex(w, site.Projects)
	}
}

func PageHandler(site app.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		site.RenderPage(w, fmt.Sprintf("%s/%s.md", app.BLOGDIR, mux.Vars(req)["pageId"]))
	}
}

func AboutSection(site app.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		site.RenderPage(w, app.ABOUTPAGE)
	}
}

func FilterByTag(site app.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		site.RenderFilterIndex(w, site.Blog, mux.Vars(req)["tagId"])
	}
}

func ViewAllPosts(site app.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		metrics.ViewIndex.Inc()
		site.RenderIndex(w, site.Blog)
	}
}
