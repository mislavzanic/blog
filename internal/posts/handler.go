package posts

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	HTMLDIR  = "html"
	CSSDIR   = "css"
	GITDIR   = "BlogPosts"
	POSTSDIR = "BlogPosts/posts"
)

type Page struct {
	Body     string
	MetaData MetaData
	Path     string
}

type Posts struct {
	Pages  []*Page
}

type MetaData struct {
	Title      string
	Date       time.Time
	Tags       []string
	TitleImage string
	Summary    string
}


func PageHandler(w http.ResponseWriter, req *http.Request) {
	pageId := mux.Vars(req)["pageId"]
	p := readBlogPost(fmt.Sprintf("%s/%s.md", POSTSDIR, pageId))

	renderFromTemplate(w, "post.html", fmt.Sprintf("%s/post.html", HTMLDIR), template.FuncMap{"markDown": markDowner}, p)
}

func FilterByTag(w http.ResponseWriter, req *http.Request) {
	tagId := mux.Vars(req)["tagId"]
	posts := findBlogPosts(tagId)

	renderFromTemplate(w, "tags.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": getUrl}, posts)
}

func ViewAllPosts(w http.ResponseWriter, req *http.Request) {
	posts := getAllPosts()
	renderFromTemplate(w, "index.html", fmt.Sprintf("%s/index.html", HTMLDIR), template.FuncMap{"toURL": getUrl}, posts)
}

func findBlogPosts(tagId string) Posts {
	posts := getAllPosts()
	p := Posts{}
	for _, post := range posts.Pages {
		for _, tag := range post.MetaData.Tags {
			if tagId == tag {
				p.Pages = append(p.Pages, post)
			}
		}
	}
	return p
}

func getAllPosts() Posts {
	files, err := os.ReadDir(POSTSDIR)

	if err != nil {
		log.Fatal(err)
	}

	posts := Posts{}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		p := readBlogPost(fmt.Sprintf("%s/%s", POSTSDIR, file.Name()))
		posts.Pages = append(posts.Pages, p)
	}

	return posts
}

func readBlogPost(path string) *Page {
	body, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	metaData := readMetadata(bytes.NewReader(body))

	return &Page{MetaData: metaData, Body: readBody(body), Path: path}
}

func readMetadata(r io.Reader) MetaData {
	metaData := make([]string, 0)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	summary := ""

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			scanner.Scan()
			summary = scanner.Text()
			break
		}
		metaData = append(metaData, line)
	}
	date, err := time.Parse("2006-01-02", metaData[1])

	if err != nil {
		log.Fatal(err)
	}

	titleImage := ""
	if len(metaData) > 3 {
		titleImage = metaData[3]
	}

	return MetaData{Title: metaData[0], Date: date, Tags: strings.Split(metaData[2], ","), TitleImage: titleImage, Summary: summary}
}

func readBody(byteArray []byte) string {
	return string(bytes.Split(byteArray, []byte("---\n"))[1])
}
