package handlers

import (
	"fmt"
	"errors"

	"codeberg.org/mislavzanic/main/internal/posts"
)

type Site struct {
	Blog posts.Posts
	About posts.Page
	Projects posts.Posts
}

func LoadSite() Site {
	return Site{
		Blog: posts.GetAllPosts(POSTSDIR),
		About: *posts.ReadBlogPost(fmt.Sprintf("%s/about.md", ABOUTDIR)),
		Projects: posts.GetAllPosts(PROJDIR),
	}
}

func (s Site) findPost(path string) (*posts.Page, error) {

	for _, post := range s.Blog.Pages {
		if post.Path == path {
			return post, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Can't find post at %s", path))
}
