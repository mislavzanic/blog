package internal

import (
	"log"
	"os"

	"codeberg.org/mislavzanic/main/internal/handlers"
	"codeberg.org/mislavzanic/main/internal/server"

	"net/http"

	"github.com/go-git/go-git/v5"
	"github.com/gorilla/mux"
)

type Blog struct {
    router *mux.Router
}

func NewBlog() *Blog {
	return &Blog{
		router: server.NewRouter(),
	}
}

func (b *Blog) Run() {
	if _, err := os.Stat(handlers.GITDIR); os.IsNotExist(err) {
		if _, err := git.PlainClone(handlers.GITDIR, false, &git.CloneOptions{
			URL:      "https://codeberg.org/mislavzanic/BlogPosts",
			Progress: os.Stdout,
		}); err != nil {
			log.Fatal(err)
		}
	}

	http.Handle("/", b.router)
	http.ListenAndServe(":8080", nil)
}
