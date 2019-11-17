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

	feature := now.Add(time.Hour * 24 * 30 * 12 * 10)
	db = db.LogMode(true)
	var post2 model.Post
	err := db.Where(model.Post{Kind: model.NativeAd}).
		Attrs("Title", "---").
		Assign(model.Post{PublishedAt: &feature}).
		FirstOrInit(&post2).Error
	if err != nil {
		panic(err)
	}
	// spew.Dump(post2)

	db.Where("author_id = ?", db.Table("authors").Where(model.Author{Name: "ymgyt"}).Select("id").Limit(1).SubQuery()).Set("gorm:auto_preload", true).First(&post2)
	spew.Dump(post2)
}
