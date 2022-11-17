package server

import (
	"net/http"

	"codeberg.org/mislavzanic/main/internal/posts"
	"codeberg.org/mislavzanic/main/internal/webhook"
	"github.com/gorilla/mux"
)


func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(posts.CSSDIR))))
	router.PathPrefix("/BlogPosts/").Handler(http.StripPrefix("/BlogPosts/", http.FileServer(http.Dir(posts.GITDIR))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(posts.JSDIR))))

	router.HandleFunc("/", posts.ViewAllPosts)
	router.HandleFunc("/about", posts.AboutSection)
	router.HandleFunc("/blog/{pageId}", posts.PageHandler)
	router.HandleFunc("/by-tag/{tagId}", posts.FilterByTag)
	router.HandleFunc("/api/wh", webhook.Webhook)
	return router
}
