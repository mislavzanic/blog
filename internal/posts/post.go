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
	metadataMap := make(map[string]string)
	metaData := MetaData{}
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			for line != "" {
				scanner.Scan()
				line = scanner.Text()
				metaData.Summary += " " + line
			}
			break
		}
		keyVal := strings.Split(line, ":")
		metadataMap[keyVal[0]] = strings.TrimSpace(keyVal[1])
	}

	if dateStr, ok := metadataMap["date"]; ok {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Fatal(err)
		}
		metaData.Date = date
	}

	if title, ok := metadataMap["title"]; ok {
		metaData.Title = title
	}

	if imgPath, ok := metadataMap["title-image"]; ok {
		metaData.TitleImage = imgPath
	}

	if tags, ok := metadataMap["tags"]; ok {
		metaData.Tags = removeEmptyStrings(strings.Split(tags, ","))
	}

	return metaData
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
