package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"codeberg.org/mislavzanic/main/internal/posts"
)


func NewRouter() *mux.Router {
	cssDir := "./css/"
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(cssDir))))
	router.HandleFunc("/blog/{pageId}", posts.PageHandler)
	return router
}
