package model

import (
	"github.com/indes/rssflow/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

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
	if config.EnableMysql {
		clientConfig := config.GetMysqlConnectingString()
		db, err := gorm.Open("mysql", clientConfig)
		if err != nil {
			panic("连接数据库失败")
		}
		return db
	}

	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("连接数据库失败")
	}
	return db
}

//EditTime timestamp
type EditTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
