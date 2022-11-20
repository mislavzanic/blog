package server

import (
	"net/http"

	"codeberg.org/mislavzanic/main/internal/handlers"
	"codeberg.org/mislavzanic/main/internal/webhook"
	"github.com/gorilla/mux"
)


func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(handlers.CSSDIR))))
	router.PathPrefix("/BlogPosts/").Handler(http.StripPrefix("/BlogPosts/", http.FileServer(http.Dir(handlers.GITDIR))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(handlers.JSDIR))))

	router.HandleFunc("/", handlers.ViewAllPosts)
	router.HandleFunc("/about", handlers.AboutSection)
	router.HandleFunc("/projects", handlers.ViewProjects)
	router.HandleFunc("/projects/{projId}", handlers.GetProject)
	router.HandleFunc("/blog/{pageId}", handlers.PageHandler)
	router.HandleFunc("/by-tag/{tagId}", handlers.FilterByTag)
	router.HandleFunc("/api/wh", webhook.Webhook)
	return router
}
