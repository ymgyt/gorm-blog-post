package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/ymgyt/gorm-blog-post/internal/lib"
	"github.com/ymgyt/gorm-blog-post/internal/model"
)

func main() {
	db := lib.Connect()

	post1 := model.Post{
		Content: "post1",
		Meta: model.Meta{
			Version: "1",
		},
		Author: model.Author{
			Name: "ymgyt",
		},
		Comments: []model.Comment{
			{Body: "comment 1"},
			{Body: "comment 2"},
		},
		Tags: []model.Tag{
			{Body: "golang"},
			{Body: "gorm"},
		},
	}

	if err := db.Save(&post1).Error; err != nil {
		panic(err)
	}

	spew.Dump(post1)
}
