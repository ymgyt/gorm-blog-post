package main

import (
	"github.com/davecgh/go-spew/spew"
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
		Base: model.Base{Meta: "metavalue"},
		Setting: model.Setting{
			Lang: "JP",
		},
	}

	desc := model.Description{Description: "about yuta..."}
	if err := db.Save(&desc).Error; err != nil {
		panic(err)
	}
	u.DescriptionID = desc.ID

	if err := db.Save(&u).Error; err != nil {
		panic(err)
	}

	var d1 model.Description
	if err := db.Model(&u).Related(&d1).Error; err != nil {
		panic(err)
	}
	spew.Dump(d1)
	// m := db.NewScope(&u).GetModelStruct()
	// spew.Dump(m)
}

func related() {
	// db := lib.Connect()

	// u := model.User{
	// 	Name:   "yuta",
	// 	MyScan: &model.MyScanner{X: 100, Y: "20"},
	// 	Profiles: []model.Profile{
	// 		{FirstName: "P1", LastName: "gopher"},
	// 	},
	// 	Base: model.Base{Meta: "metavalue"},
	// 	Setting: model.Setting{
	// 		Lang: "JP",
	// 	},
	// }

	// desc := model.Description{Description: "about yuta..."}
	// if err := db.Save(&desc).Error; err != nil {
	// 	panic(err)
	// }
	// u.DescriptionID = desc.ID

	// if err := db.Save(&u).Error; err != nil {
	// 	panic(err)
	// }

	// var d1 model.Description
	// if err := db.Model(&u).Related(&d1).Error; err != nil {
	// 	panic(err)
	// }
	// spew.Dump(d1)

}
