package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type User struct {
	ID         uint64
	Name       string
	HasDefault string `gorm:"DEFAULT"`
	MyScan     *MyScanner
	Profiles   []Profile
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MyScanner struct {
	X int
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
