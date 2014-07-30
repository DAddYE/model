package model_test

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/daddye/model"
	"github.com/daddye/model/sql"
	_ "github.com/lib/pq"
)

type Feed struct {
	Id               int    `sql:"id"`
	SoruceFileFormat string `sql:"source_file_format"`
	SourceType       string `sql:"source_type"`
	Ftp              Ftp
	model.Model
}

func (f *Feed) First(conditions string, args ...interface{}) error {
	return query.First(&f.Model, conditions, args...)
}

func (f *Feed) Find(conditions string, args ...interface{}) (*model.Iter, error) {
	return query.Find(&f.Model, conditions, args...)
}

type Ftp struct {
	FileName string `sql:"source_ftp_file_name"`
	Username string `sql:"source_ftp_username"`
}

func NewFeed(db *sql.DB) (f *Feed) {
	f = new(Feed)
	f.Interface = db
	f.Table = "feeds"
	model.Set(f, "sql")
	return
}

func ExampleMulti() {
	db, err := sql.Open("postgres", "user=... dbname=... sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	feed := NewFeed(db)
	rows, err := feed.Find("WHERE source_type = 'ftp' LIMIT $1", 2)
	if err != nil {
		fmt.Print(err)
		return
	}
	for rows.Next() {
		fmt.Printf("Id: %d, FileName: %s\n", feed.Id, feed.Ftp.Username)
	}
}

func ExampleSingle() {
	db, err := sql.Open("postgres", "user=... dbname=... sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	feed := NewFeed(db)
	err = feed.First("WHERE source_type = 'ftp'")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("Id: %d, FileName: %s\n", feed.Id, feed.Ftp.Username)
}
