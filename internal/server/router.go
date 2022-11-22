package server

import (
	"net/http"

	"codeberg.org/mislavzanic/main/internal/handlers"
	"github.com/gorilla/mux"
)


func NewRouter(site handlers.Site) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(handlers.CSSDIR))))
	router.PathPrefix("/blog/").Handler(http.StripPrefix("/blog/", http.FileServer(http.Dir(handlers.GITDIR))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(handlers.JSDIR))))

	router.HandleFunc("/", site.ViewAllPosts)
	router.HandleFunc("/about", site.AboutSection)
	router.HandleFunc("/projects", site.ViewProjects)
	router.HandleFunc("/posts/{pageId}", site.PageHandler)
	router.HandleFunc("/by-tag/{tagId}", site.FilterByTag)
	return router
}
