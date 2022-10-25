package posts

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

type Page struct {
	Body     string
	MetaData MetaData
}

type MetaData struct {
	Title string
	Date  time.Time
	Tags  []string
}

func markDowner(args ...interface{}) template.HTML {
	return template.HTML(blackfriday.Run([]byte(fmt.Sprintf("%s", args...))))
}

func PageHandler(w http.ResponseWriter, req *http.Request) {
	pageId := mux.Vars(req)["pageId"]

	p := readBlogPost(fmt.Sprintf("./posts/%s.md", pageId))


	html_template, _ := ioutil.ReadFile("html/post.html")
	tmpl := template.Must(template.New("test.html").Funcs(template.FuncMap{"markDown": markDowner}).Parse(string(html_template)))
	err := tmpl.ExecuteTemplate(w, "test.html", p)

	if err != nil {
		log.Fatal(err)
	}
}

func readBlogPost(path string) *Page {
	body, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	metaData := readMetadata(bytes.NewReader(body))

	return &Page{MetaData: metaData, Body: readBody(body)}
}

func readMetadata(r io.Reader) MetaData {
	metaData := make([]string, 0)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			break
		}
		metaData = append(metaData, line)
	}
	date, err := time.Parse("2006-01-02", metaData[1])

	if err != nil {
		log.Fatal(err)
	}

	return MetaData{Title: metaData[0], Date: date, Tags: strings.Split(metaData[2], ",")}
}

func readBody(byteArray []byte) string {
	return string(bytes.Split(byteArray, []byte("---\n"))[1])
}
