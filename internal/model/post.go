package model

import (
	"fmt"
	"strconv"
	"time"
)

type Meta struct {
	ID           uint
	Version      string
	ResourceID   uint
	ResourceType string
}

type Author struct {
	ID   uint
	Name string
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

type PostKind int

const (
	Undefined PostKind = iota
	Normal
	NativeAd
)

func (pk *PostKind) Scan(src interface{}) error {
	parse := func(n int) {
		switch PostKind(n) {
		case Normal:
			*pk = Normal
		case NativeAd:
			*pk = NativeAd
		default:
			*pk = Undefined
		}
	}

	switch v := src.(type) {
	case int:
		parse(v)
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to scan PostKind: %w", err)
		}
		parse(int(n))
	}

	return nil
}

type Style struct {
	Font  string
	Theme string
}

type Content struct {
	ID      uint
	PostRef uint
	Body    string
}

type Post struct {
	ID    uint
	Kind  PostKind
	Title string

	// belong_to
	AuthorID uint
	Author   Author

	// has_one
	Content Content `gorm:"foreignkey:PostRef"`

	// has_one polymorphic
	Meta Meta `gorm:"polymorphic:Resource;polymorphic_value:post"`

	// has_many
	Comments []Comment

	// many2many
	Tags []Tag `gorm:"many2many:post_tags"`

	// tags
	IgnoreMe int `gorm:"-"`

	Style Style `gorm:"EMBEDDED; EMBEDDED_PREFIX:post_"`

	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// func (p Post) String() string {
// 	return "aaaa"
// }
