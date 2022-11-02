package main

import (
	"codeberg.org/mislavzanic/main/internal"
)

func main() {
	b := internal.NewBlog()
	b.Run()
}
