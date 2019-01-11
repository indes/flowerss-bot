package model

import (
	"github.com/SlyMarbo/rss"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
)

var ()

func init() {
	db := getConnect()
	defer db.Close()
	db.LogMode(true)
	if !db.HasTable(&Source{}) {
		db.CreateTable(&Source{})
	}

	if !db.HasTable(&Subscribe{}) {
		db.CreateTable(&Subscribe{})
	}

	if !db.HasTable(&Content{}) {
		db.CreateTable(&Content{})
	}
}

type Source struct {
	ID         uint `gorm:"primary_key";"AUTO_INCREMENT"`
	Link       string
	Title      string
	ErrorCount uint
	Content    []Content
	EditTime
}

type Subscribe struct {
	UserID   int64
	SourceID uint
	EditTime
}

type Content struct {
	SourceID     uint
	HashID       string `gorm:"primary_key"`
	RawLink      string
	TelegraphUrl string
	EditTime
}

type EditTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func RegistFeed(userID int64, sourceID uint) error {
	var subscribe Subscribe
	db := getConnect()
	defer db.Close()

	if err := db.Where("user_id=? and source_id=?", userID, sourceID).Find(&subscribe).Error; err != nil {
		if err.Error() == "record not found" {
			subscribe.UserID = userID
			subscribe.SourceID = sourceID
			err := db.Create(&subscribe).Error
			if err == nil {
				return nil
			}
		}
		return err
	}
	return nil
}

func FindOrNewSourceByUrl(url string) (*Source, error) {
	var source Source
	db := getConnect()
	defer db.Close()

	if err := db.Where("link=?", url).Find(&source).Error; err != nil {
		if err.Error() == "record not found" {
			source.Link = url

			// parsing rss
			feed, err := rss.Fetch(url)
			if err != nil {
				log.Println("Unable to make request: ", err)
				return nil, err
			}

			source.Title = feed.Title
			source.ErrorCount = 0

			db.Create(&source)
			return &source, nil
		}
		return nil, err
	}

	return &source, nil
}

func getConnect() *gorm.DB {
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("连接数据库失败")
	}
	return db
}
