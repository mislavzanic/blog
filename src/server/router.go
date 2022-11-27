package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mislavzanic/blog/src/app"
	"github.com/mislavzanic/blog/src/handlers"
)


func NewRouter(site app.Site) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(handlers.CSSDIR))))
	router.PathPrefix("/post/").Handler(http.StripPrefix("/post/", http.FileServer(http.Dir("webContent"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(handlers.JSDIR))))

	router.HandleFunc("/", handlers.ViewAllPosts(site))
	router.HandleFunc("/about", handlers.AboutSection(site))
	router.HandleFunc("/projects", handlers.ViewProjects(site))
	router.HandleFunc("/blog/{pageId}", handlers.PageHandler(site))
	router.HandleFunc("/by-tag/{tagId}", handlers.FilterByTag(site))

	http.Handle("/metrics", promhttp.Handler())

	return router
}
