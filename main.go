package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/russross/blackfriday/v2"
	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  string
}

func hello(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadFile("./posts/test.md")
	p := &Page{Title: "A Test Demo", Body: string(body)}

	html_template, _ := ioutil.ReadFile("html/post.html")

	tmpl := template.Must(template.New("test.html").Funcs(template.FuncMap{"markDown": markDowner}).Parse(string(html_template)))
	err := tmpl.ExecuteTemplate(w, "test.html", p)
	if err != nil {
		fmt.Println(err)
	}
}

func markDowner(args ...interface{}) template.HTML {
	return template.HTML(blackfriday.Run([]byte(fmt.Sprintf("%s", args...))))
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	r := mux.NewRouter()
	cssHandler := http.FileServer(http.Dir("./css/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	r.HandleFunc("/", hello)
	http.Handle("/", r)
	http.ListenAndServe(":8070", nil)
}
