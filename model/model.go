package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

type EditTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
