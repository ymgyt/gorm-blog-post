package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ymgyt/gorm-blog-post/internal/lib"
	"github.com/ymgyt/gorm-blog-post/internal/model"
)

func main() {
	db := lib.Connect()
	db = db.LogMode(false)
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

	if err := db.Save(&post1).Error; err != nil {
		panic(err)
	}

	db = db.LogMode(true)
	// var author model.Author
	// err := db.Model(&post1).Association("Author").Find(&author).Error
	// if err != nil {
	// 	panic(err)
	// }

	if err := db.Model(&post1).Association("Tags").Append(model.Tag{ID: 3, Body: "new"}).Error; err != nil {
		panic(err)
	}
	spew.Dump(post1)
}
