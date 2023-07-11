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
	router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir(app.CSSDIR))))
	router.PathPrefix("/post/").Handler(http.StripPrefix("/post/", http.FileServer(http.Dir("webContent"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(app.JSDIR))))
	router.PathPrefix("/static/img/").Handler(http.StripPrefix("/static/img/", http.FileServer(http.Dir(app.IMGDIR))))

	router.HandleFunc("/", handlers.AboutSection(site))
	router.HandleFunc("/blog", handlers.ViewAllPosts(site))
	router.HandleFunc("/blog/{pageId}", handlers.PageHandler(site))
	router.HandleFunc("/by-tag/{tagId}", handlers.FilterByTag(site))
	router.HandleFunc("/projects", handlers.ViewProjects(site))

	http.Handle("/metrics", promhttp.Handler())

	return router
}
