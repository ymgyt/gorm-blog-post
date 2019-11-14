package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ymgyt/gorm-blog-post/internal/lib"
	"github.com/ymgyt/gorm-blog-post/internal/model"
)

func main() {
	db := lib.Connect()
	now := time.Now()

	post1 := model.Post{
		Kind: model.Normal,
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
		PublishedAt: &now,
	}

	db = db.LogMode(false)
	if err := db.Save(&post1).Error; err != nil {
		panic(err)
	}

	db = db.LogMode(true)
	var post2 model.Post
	if err := db.
		Where(model.Post{ID: post1.ID}).
		Preload("Author").
		Preload("Content").
		Preload("Meta").
		Preload("Comments").
		Preload("Tags").
		Find(&post2).Error; err != nil {
		panic(err)
	}
	spew.Dump(post2)
}
