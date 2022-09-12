package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"moul.io/zapgorm"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/pkg/client"
)

var db *gorm.DB
var httpClient *client.HttpClient // TODO: 将网络拉取逻辑从 model 包移除

// InitDB init db object
func InitDB() {
	connectDB()
	configDB()
	updateTable()
	initHttpClient()
}

func initHttpClient() {
	clientOpts := []client.HttpClientOption{
		client.WithTimeout(10 * time.Second),
	}
	if config.Socks5 != "" {
		clientOpts = append(clientOpts, client.WithProxyURL(fmt.Sprintf("socks5://%s", config.Socks5)))
	}

	if config.UserAgent != "" {
		clientOpts = append(clientOpts, client.WithUserAgent(config.UserAgent))
	}
	httpClient = client.NewHttpClient(clientOpts...)
}

func configDB() {
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	db.LogMode(config.DBLogMode)
	db.SetLogger(zapgorm.New(log.Logger.WithOptions(zap.AddCallerSkip(7))))
}

func updateTable() {
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
		zap.S().Fatalf("connect db failed, err: %+v", err)
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

// EditTime timestamp
type EditTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
