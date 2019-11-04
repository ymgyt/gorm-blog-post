package main

import (
	"os"

	"github.com/ymgyt/gorm-blog-post/internal/lib"

)

func main() {
	if err := lib.GenerateDBConf(os.Stdout); err != nil {
		panic(err)
	}
}
