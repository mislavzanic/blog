package posts

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

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

func FindBlogPosts(tagId, dir string) Posts {
	posts := GetAllPosts(dir)
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

func GetAllPosts(dir string) Posts {
	files, err := os.ReadDir(dir)
	posts := Posts{}

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		p := ReadBlogPost(fmt.Sprintf("%s/%s", dir, file.Name()))
		posts.Pages = append(posts.Pages, p)
	}

	sort.Slice(posts.Pages, func(i, j int) bool {
		return posts.Pages[i].MetaData.Date.After(posts.Pages[j].MetaData.Date)
	})

	return posts
}

func ReadBlogPost(path string) *Page {
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
