package model

import "time"

type Meta struct {
	ID      uint
	Version string
}

func (m Meta) TableName() string { return "metas" }

type Author struct {
	ID     uint
	Name   string
	PostID uint
}

type Comment struct {
	ID     uint
	Body   string
	PostID uint
}

type Tag struct {
	ID   uint
	Body string
}

type Post struct {
	ID      uint
	Content string

	// belong_to
	MetaID uint
	Meta   Meta

	// has_one
	Author Author

	// has_many
	Comments []Comment

	// many2many
	Tags []Tag `gorm:"many2many:post_tags"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
