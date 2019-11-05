package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Base struct {
	Meta string
}

type Setting struct {
	ID     uint64
	UserID uint64
	Lang   string
}

type User struct {
	ID         uint64
	Name       string
	HasDefault string     `gorm:"DEFAULT"`
	MyScan     *MyScanner `gorm:"MyScannerOuterKey:MyScannerOuterValue"`
	Profiles   []Profile

	// Memo
	// DBのuser_meta fieldを対象にする
	// prefix + embedded struct field
	Base Base `gorm:"EMBEDDED;EMBEDDED_PREFIX:user_"`

	Setting Setting

	// Anonymous struct {
	// 	AnonymousField string
	// }

	// Memo
	// time.TimeはField.StructField.IsNormal = trueが入る
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MEMO: MyScan fieldに定義したtagとmergeされる。
// User.MyScan.Tag
// TagSettings: (map[string]string) (len=2) {
// 	(string) (len=17) "MYSCANNEROUTERKEY": (string) (len=19) "MyScannerOuterValue",
// 	(string) (len=17) "MYSCANNERINNERKEY": (string) (len=19) "MyScannerInnerValue"
// },
type MyScanner struct {
	X int `gorm:"MyScannerInnerKey:MyScannerInnerValue"`
	Y string
}

func (s *MyScanner) Scan(src interface{}) error {
	fmt.Printf("MyScanner.Scan(%v)\n", src)
	return nil
}

func (s *MyScanner) Value() (driver.Value, error) {
	return driver.Value(fmt.Sprintf("x:%d y:%s", s.X, s.Y)), nil
}

type Profile struct {
	ID        uint64
	UserID    uint64
	FirstName string
	LastName  string
}
