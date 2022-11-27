package posts

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type void struct{}
type set map[string]void

type Page struct {
	Body     string
	Summary  string
	ReadTime string
	MetaData MetaData
	Path     string
}

type Posts struct {
	Pages   []*Page
	Uri     string
	AllTags []string
}

type MetaData struct {
	Title      string    `yaml:"title"`
	Date       time.Time `yaml:"date"`
	Tags       []string  `yaml:"tags"`
	TitleImage string    `yaml:"title-image"`
	Link       string    `yaml:"link"`
	Latex      bool      `yaml:"latex"`
}

func (p Posts) FindBlogPosts(tagId string) Posts {
	index := Posts{}
	index.AllTags = p.AllTags
	index.Uri = p.Uri
	for _, post := range p.Pages {
		for _, tag := range post.MetaData.Tags {
			if tagId == tag {
				index.Pages = append(index.Pages, post)
			}
		}
	}
	return index
}

func GetAllPosts(dir string) Posts {
	files, err := os.ReadDir(dir)
	posts := Posts{}
	tags := make(map[string]struct{})
	posts.Uri = strings.Split(dir, "/")[1]

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		p := ReadBlogPost(fmt.Sprintf("%s/%s", dir, file.Name()))
		for _, tag := range p.MetaData.Tags {
			tags[tag] = void{};
		}
		posts.Pages = append(posts.Pages, p)
	}

	sort.Slice(posts.Pages, func(i, j int) bool {
		return posts.Pages[i].MetaData.Date.After(posts.Pages[j].MetaData.Date)
	})

	posts.AllTags = make([]string, 0, len(tags))
	for tag := range tags {
		posts.AllTags = append(posts.AllTags, tag)
	}

	sort.Slice(posts.AllTags, func(i, j int) bool {
		return posts.AllTags[i] < posts.AllTags[j]
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

	return &Page{
		MetaData: metaData,
		Body: postBody,
		Summary: summary,
		Path: path,
		ReadTime: strconv.Itoa((len(strings.Fields(postBody)) / 300) + 1),
	}
}

func readMetadata(body []byte) MetaData {
	yamlPart := bytes.Split(body, []byte("---\n"))[1]

	var metadata MetaData
	if err := yaml.Unmarshal(yamlPart, &metadata); err != nil {
		log.Fatal(err)
	}

	return metadata
}

func readBody(byteArray []byte) (string, string) {
	body := string(bytes.Split(byteArray, []byte("---\n"))[2])
	summary := strings.Split(body, "\n\n")[0]
	return body, summary
}
