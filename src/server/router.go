package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mislavzanic/blog/src/app"
	"github.com/mislavzanic/blog/src/handlers"
)


func NewRouter(site app.Site) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(handlers.CSSDIR))))
	router.PathPrefix("/blog/").Handler(http.StripPrefix("/blog/", http.FileServer(http.Dir("blog"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(handlers.JSDIR))))

	router.HandleFunc("/", handlers.ViewAllPosts(site))
	router.HandleFunc("/about", handlers.AboutSection(site))
	router.HandleFunc("/projects", handlers.ViewProjects(site))
	router.HandleFunc("/posts/{pageId}", handlers.PageHandler(site))
	router.HandleFunc("/by-tag/{tagId}", handlers.FilterByTag(site))
	return router
}
