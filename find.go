package main

import (
	"fmt"
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
			Reviews: []model.Review{
				{Body: "review 1"},
			},
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

	db = db.LogMode(true)
	if err := db.Save(&post1).Error; err != nil {
		panic(err)
	}

	db = db.LogMode(true)
	var post2 model.Post
	db = db.
		Where(model.Post{ID: post1.ID}).
		Preload("Author").
		Preload("Author.Reviews").
		// Preload("Content").
		// Preload("Meta").
		// Preload("Comments").
		// Preload("Tags").
		// Set("gorm:auto_preload", true).
		Set("gorm:query_option", "FOR UPDATE").
		Find(&post2)
	if db.Error != nil {
		panic(db.Error)
	}
	fmt.Println(db.RowsAffected)
	spew.Dump(post2)
}
