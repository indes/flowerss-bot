package model

import (
	"fmt"
	"net/url"

	"encoding/csv"
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// const sqliteFulltextTable = "CREATE VIRTUAL TABLE feed_fts USING FTS5(name,bizid,description)"
const (
	wechatAccountSource = "https://raw.githubusercontent.com/hellodword/wechat-feeds/main/list.csv"
	wechatHost          = "mp.weixin.qq.com"
	WechatSubUrl        = "https://github.com/hellodword/wechat-feeds/raw/feeds/%s.xml"
)

type Feed struct {
	ID          int    `gorm:"primaryKey,autoIncrement"`
	Name        string `gorm:"class:FULLTEXT"`
	Bizid       string `gorm:"uniqueIndex"`
	Description string
}

// ProcessWechatURL return wechat-feed sub URL. if it's valid wehchat url, else return origin str
func ProcessWechatURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err == nil {
		if u.Host == "mp.weixin.qq.com" {
			q := u.Query()
			bizs, ok := q["__biz"]
			if ok {
				biz := bizs[0]
				newURL := fmt.Sprintf(WechatSubUrl, biz)
				return newURL
			}
		}
	}
	return urlStr
}

func loadData() []Feed {
	resp, err := http.Get(wechatAccountSource)
	if err != nil {
		log.Fatal("Get wechat account from origin source failed. ", err)
		return nil
	}

	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse csv file. ", err)
	}
	var feeds []Feed
	for _, record := range records[1:] {
		feed := Feed{
			// ID:          idx,
			Name:        record[0],
			Bizid:       record[1],
			Description: record[2],
		}
		feeds = append(feeds, feed)
	}
	return feeds
}

// LoadWechatAccounts 导入所有的可以订阅的公众号信息，明天更新一次
func LoadWechatAccounts() {

	feeds := loadData()
	db.AutoMigrate(&Feed{})
	db.DropTable(&Feed{})
	db.CreateTable(&Feed{})

	tx := db.Begin()
	for _, feed := range feeds {
		tx.Create(&feed)
	}
	tx.Commit()
}

func SearchWechatAccounts(keyword string, offset, limit int) ([]Feed, int64) {
	keyword = "%" + keyword + "%"

	var total int64
	q := db.Model(&Feed{}).Where("name like ? or bizid like ? or description like ? ", keyword, keyword, keyword)
	q.Count(&total)
	rows, err := q.Offset(offset).Limit(limit).Rows()
	if err != nil {
		log.Fatal("Unable to search query with keyword: "+keyword+".\t", err)
	}
	defer rows.Close()

	var feeds []Feed
	var feed Feed
	for rows.Next() {
		db.ScanRows(rows, &feed)
		feeds = append(feeds, feed)
	}
	return feeds, total
}
