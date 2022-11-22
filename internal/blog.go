package internal

import (
	"codeberg.org/mislavzanic/main/internal/handlers"
	"codeberg.org/mislavzanic/main/internal/server"

	"net/http"

	"github.com/gorilla/mux"
)

type Blog struct {
    router *mux.Router
}

func NewBlog() *Blog {
	site := handlers.LoadSite()

	return &Blog{
		router: server.NewRouter(site),
	}
}

func (b *Blog) Run() {
	http.Handle("/", b.router)
	http.ListenAndServe(":8080", nil)
}
