package main

import (
	"github.com/ymgyt/gorm-blog-post/internal/lib"
	"github.com/ymgyt/gorm-blog-post/internal/model"
)

func main() {
	db := lib.Connect()

	u := model.User{
		Name:   "yuta",
		MyScan: &model.MyScanner{X: 100, Y: "20"},
		Profiles: []model.Profile{
			{FirstName: "P1", LastName: "gopher"},
		},
	}

	if err := db.Save(&u).Error; err != nil {
		panic(err)
	}
}
