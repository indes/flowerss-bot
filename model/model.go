package model

import (
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.uber.org/zap"
	"moul.io/zapgorm"

	"time"
)

var db *gorm.DB

// InitDB init db object
func InitDB() {
	connectDB()
	configDB()
	updateTable()
}

func configDB(){
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	db.LogMode(config.DBLogMode)
	db.SetLogger(zapgorm.New(log.Logger.WithOptions(zap.AddCallerSkip(7))))
}

func updateTable(){
	createOrUpdateTable(&Subscribe{})
	createOrUpdateTable(&User{})
	createOrUpdateTable(&Source{})
	createOrUpdateTable(&Option{})
	createOrUpdateTable(&Content{})
}

// connectDB connect to db
func connectDB() {
	if config.RunMode == config.TestMode {
		return
	}

	var err error
	if config.EnableMysql {
		db, err = gorm.Open("mysql", config.Mysql.GetMysqlConnectingString())
	} else {
		db, err = gorm.Open("sqlite3", config.SQLitePath)
	}
	if err != nil {
		log.Fatal(err.Error())
	}
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
