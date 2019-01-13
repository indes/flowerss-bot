package model

import (
	"github.com/SlyMarbo/rss"
	tgp "github.com/indes/rssflow/tgraph"
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
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}
}

func getConnect() *gorm.DB {
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("连接数据库失败")
	}
	return db
}

type User struct {
	ID     int64    `gorm:"primary_key"`
	Source []Source `gorm:"many2many:subscribes;"`
	EditTime
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
	ID       uint `gorm:"primary_key";"AUTO_INCREMENT"`
	UserID   int64
	SourceID uint
	EditTime
}

type Content struct {
	SourceID     uint
	HashID       string `gorm:"primary_key"`
	RawID        string
	RawLink      string
	Title        string
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

			// Get contents and insert
			items := feed.Items
			source.appendContents(items)
			tgp.PublishItems(items)
			db.Create(&source)
			return &source, nil
		}
		return nil, err
	}

	return &source, nil
}

func (s *Source) appendContents(items []*rss.Item) error {
	for _, item := range items {
		var c = Content{
			Title:        item.Title,
			SourceID:     s.ID,
			RawID:        item.ID,
			HashID:       genHashID(s, item.ID),
			RawLink:      item.Link,
			TelegraphUrl: item.Link,
		}

		s.Content = append(s.Content, c)
	}

	return nil
}

func GetSubscribeByUserID(userID int64) []Source {
	db := getConnect()
	defer db.Close()
	user := FindOrInitUser(userID)
	return user.Source
}

func FindOrInitUser(userID int64) *User {
	db := getConnect()
	defer db.Close()
	var user User
	//db.FirstOrInit(User{ID: userID}, &user)
	db.Where(User{ID: userID}).FirstOrCreate(&user)
	return &user
}
