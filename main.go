package main

import (
	"net/http"
	"codeberg.org/mislavzanic/main/internal/server"
)

func main() {
	r := server.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
