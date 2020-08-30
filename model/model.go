package model

import (
	"github.com/indes/flowerss-bot/config"
	"github.com/jinzhu/gorm"

	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

// ConnectDB connect to db and update table
func ConnectDB() {
	var err error

	if config.EnableMysql {
		db, err = gorm.Open("mysql", config.Mysql.GetMysqlConnectingString())
	} else {
		db, err = gorm.Open("sqlite3", config.SQLitePath)
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	db.LogMode(true)

	createOrUpdateTable(&Subscribe{})
	createOrUpdateTable(&User{})
	createOrUpdateTable(&Source{})
	createOrUpdateTable(&Option{})
	createOrUpdateTable(&Content{})
}

// Disconnect disconnects from the database.
func Disconnect() {
	db.Close()
}

// createOrUpdateTable create table or Migrate table
func createOrUpdateTable(model interface{}) {
	if !db.HasTable(model) {
		db.CreateTable(model)
	} else {
		db.AutoMigrate(model)
	}
}

//EditTime timestamp
type EditTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
