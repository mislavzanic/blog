package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/russross/blackfriday/v2"
)

func hello(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadFile("./README.md")
	output := blackfriday.Run(body)
	fmt.Fprintf(w, "%v", template.HTML(output))
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8093", nil)
}
