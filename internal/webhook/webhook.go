package webhook

import (
	"net/http"
	"os"
	"log"

	"github.com/go-git/go-git/v5"
	"codeberg.org/mislavzanic/main/internal/handlers"
)


func Webhook(w http.ResponseWriter, req *http.Request) {
	if _, err := os.Stat(handlers.GITDIR); !os.IsNotExist(err) {
		r, err := git.PlainOpen(handlers.GITDIR)

		if err != nil {
			log.Fatal(err)
		}

		w, err := r.Worktree()

		if err != nil {
			log.Fatal(err)
		}

		if err := w.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
			log.Fatal(err)
		}
	}
}
