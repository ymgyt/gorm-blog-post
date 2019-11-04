package lib

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"text/template"

	"github.com/ymgyt/gorm"
	_ "github.com/ymgyt/gorm/dialects/mysql"
)

func Connect() *gorm.DB {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	val := url.Values{}
	val.Add("charset", "utf8mb4")
	val.Add("parseTime", "True")
	val.Add("loc", "Asia/Tokyo")
	dsn := fmt.Sprintf("%s?%s", conn, val.Encode())

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	return db
}

func GenerateDBConf(w io.Writer) error {
	const src = `# goose
development:
  driver: mysql
  open: {{.User}}:{{.Pass}}@tcp({{.Host}}:{{.Port}})/{{.Name}}
`
	t, err := template.New("dbconfig").Parse(src)
	if err != nil {
		return err
	}
	return t.Execute(w, struct {
		User, Pass, Host, Port, Name string
	}{
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Name: os.Getenv("DB_NAME"),
	})
}
