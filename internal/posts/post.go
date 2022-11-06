package posts

import (
	"fmt"
	"bufio"
	"sort"
	"time"
	"bytes"
	"io"
	"os"
	"io/ioutil"
	"log"
	"strings"
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

	date, err := time.Parse("2006-01-02", strings.TrimSpace(strings.Split(metaData[1], ":")[1]))

	if err != nil {
		log.Fatal(err)
	}

	titleImage := ""
	if len(metaData) > 3 {
		titleImage = strings.TrimSpace(strings.Split(metaData[3], ":")[1])
	}

	tags := strings.Split(strings.TrimSpace(strings.Split(metaData[2], ":")[1]), ",")
	tags = removeEmptyStrings(tags)
	return MetaData{Title: strings.TrimSpace(strings.Split(metaData[0], ":")[1]), Date: date, Tags: tags,TitleImage: titleImage, Summary: summary}
}

func readBody(byteArray []byte) string {
	return string(bytes.Split(byteArray, []byte("---\n"))[1])
}

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
