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

	db = db.LogMode(false)
	if err := db.Save(&post1).Error; err != nil {
		panic(err)
	}

	db = db.LogMode(true)
	// only Post field
	post1.Kind = model.NativeAd
	if err := db.Set("gorm:association_save_reference", false).Save(&post1).Error; err != nil {
		panic(err)
	}
	spew.Dump(post1)
}
