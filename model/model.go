package model

import (
	"github.com/indes/flowerss-bot/config"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
)

func init() {
	db := getConnect()
	defer db.Close()

	db.LogMode(true)

	createOrUpdateTable(db, &Subscribe{})
	createOrUpdateTable(db, &User{})
	createOrUpdateTable(db, &Source{})
	createOrUpdateTable(db, &Option{})
	createOrUpdateTable(db, &Content{})

}

// createOrUpdateTable create table or Migrate table
func createOrUpdateTable(db *gorm.DB, model interface{}) {
	if !db.HasTable(model) {
		db.CreateTable(model)
	} else {
		db.AutoMigrate(model)
	}
}

func getConnect() *gorm.DB {
	if config.EnableMysql {
		clientConfig := config.Mysql.GetMysqlConnectingString()
		db, err := gorm.Open("mysql", clientConfig)
		if err != nil {
			panic("连接MySQL数据库失败")
		}
		return db
	} else {
		db, err := gorm.Open("sqlite3", config.SQLitePath)
		if err != nil {
			log.Println(err.Error())
			panic("连接SQLite数据库失败")
		}
		return db
	}
}

//EditTime timestamp
type EditTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
