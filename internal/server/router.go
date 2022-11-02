package server

import (
	"net/http"

	"codeberg.org/mislavzanic/main/internal/posts"
	"codeberg.org/mislavzanic/main/internal/webhook"
	"github.com/gorilla/mux"
)


func NewRouter() *mux.Router {
	cssDir := posts.CSSDIR
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(cssDir))))
	router.PathPrefix("/BlogPosts/").Handler(http.StripPrefix("/BlogPosts/", http.FileServer(http.Dir(posts.GITDIR))))
	router.HandleFunc("/", posts.ViewAllPosts)
	router.HandleFunc("/blog/{pageId}", posts.PageHandler)
	router.HandleFunc("/by-tag/{tagId}", posts.FilterByTag)
	router.HandleFunc("/api/wh", webhook.Webhook)
	return router
}
