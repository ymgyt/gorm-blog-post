package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/ymgyt/gorm-blog-post/internal/lib"
	"github.com/ymgyt/gorm-blog-post/internal/model"
)

func main() {
	db := lib.Connect()

	post1 := model.Post{
		Content: model.Content{
			Body: "post 1 content",
		},
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

		IgnoreMe: 999,
		Style: model.Style{
			Font:  "ricty",
			Theme: "Solarized dark",
		},
	}

	if err := db.Save(&post1).Error; err != nil {
		panic(err)
	}
	spew.Dump(db.NewScope(&post1).Fields())

}
