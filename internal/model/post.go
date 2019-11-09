package model

import (
	"fmt"
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
	n, ok := src.(int)
	if !ok {
		return fmt.Errorf("invalid src %v", src)
	}
	switch v := PostKind(n); v {
	case Normal:
		*pk = Normal
	case NativeAd:
		*pk = NativeAd
	default:
		*pk = Undefined
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

	CreatedAt time.Time
	UpdatedAt time.Time
}
