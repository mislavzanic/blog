package posts

import (
	"fmt"
	"strings"
	"bytes"
	"sort"
	"time"
	"io/ioutil"
	"log"
	"os"
	"gopkg.in/yaml.v2"
)

type Page struct {
	Body     string
	Summary  string
	MetaData MetaData
	Path     string
}

type Posts struct {
	Pages  []*Page
}

type MetaData struct {
	Title      string    `yaml:"title"`
	Date       time.Time `yaml:"date"`
	Tags       []string  `yaml:"tags"`
	TitleImage string    `yaml:"title-image"`
}

func findBlogPosts(tagId string) Posts {
	posts := getAllPosts(POSTSDIR)
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

func getAllPosts(dir string) Posts {
	files, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	posts := Posts{}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		p := readBlogPost(fmt.Sprintf("%s/%s", dir, file.Name()))
		posts.Pages = append(posts.Pages, p)
	}

	sort.Slice(posts.Pages, func(i, j int) bool {
		return posts.Pages[i].MetaData.Date.After(posts.Pages[j].MetaData.Date)
	})

	return posts
}

func readBlogPost(path string) *Page {
	body, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	metaData := readMetadata(body)
	postBody, summary := readBody(body)

	return &Page{MetaData: metaData, Body: postBody, Summary: summary, Path: path}
}

func readMetadata(body []byte) MetaData {
	yamlPart := bytes.Split(body, []byte("---\n"))[0]

	var metadata MetaData
	if err := yaml.Unmarshal(yamlPart, &metadata); err != nil {
		log.Fatal(err)
	}

	return metadata
}

func readBody(byteArray []byte) (string, string) {
	body := string(bytes.Split(byteArray, []byte("---\n"))[1])
	summary := strings.Split(body, "\n\n")[0]
	return body, summary
}
