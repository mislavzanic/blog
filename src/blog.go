package blog

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mislavzanic/blog/src/app"
	"github.com/mislavzanic/blog/src/server"
)

type Blog struct {
    router *mux.Router
}

func NewBlog() *Blog {
	site := app.LoadSite()

	return &Blog{
		router: server.NewRouter(site),
	}
}

func (b *Blog) Run() {
	http.Handle("/", b.router)
	http.ListenAndServe(":8080", nil)
}
