package model

import (
	"errors"
	"fmt"
	"github.com/SlyMarbo/rss"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/util"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unicode"
)

type Source struct {
	ID         uint `gorm:"primary_key;AUTO_INCREMENT"`
	Link       string
	Title      string
	ErrorCount uint
	Content    []Content
	EditTime
}

func (s *Source) appendContents(items []*rss.Item) error {
	db := getConnect()
	defer db.Close()
	for _, item := range items {
		c, _ := getContentByFeedItem(s, item)
		s.Content = append(s.Content, c)
	}
	// 开启task更新
	s.ErrorCount = 0
	db.Save(&s)
	return nil
}

func GetSourceByUrl(url string) (*Source, error) {
	var source Source
	db := getConnect()
	defer db.Close()
	if err := db.Where("link=?", url).Find(&source).Error; err != nil {
		return nil, err
	}
	return &source, nil
}

func fetchFunc(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	if config.UserAgent != "" {
		req.Header.Set("User-Agent", config.UserAgent)
	} else {
		req.Header.Set("User-Agent", "flowerss/2.0")
	}

	resp, err = util.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []byte
	if data, err = ioutil.ReadAll(resp.Body); err != nil {

		return nil, err
	}

	resp.Body = ioutil.NopCloser(strings.NewReader(strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, string(data))))
	return
}

func FindOrNewSourceByUrl(url string) (*Source, error) {
	var source Source
	db := getConnect()
	defer db.Close()

	if err := db.Where("link=?", url).Find(&source).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			source.Link = url

			// parsing task
			feed, err := rss.FetchByFunc(fetchFunc, url)

			if err != nil {
				return nil, fmt.Errorf("Feed 抓取错误 %v", err)
			}

			source.Title = feed.Title
			// 避免task更新
			source.ErrorCount = config.ErrorThreshold + 1

			// Get contents and insert
			items := feed.Items
			db.Create(&source)
			go source.appendContents(items)
			return &source, nil
		}
		return nil, err
	}

	return &source, nil
}

func GetSources() []Source {
	var sources []Source

	db := getConnect()
	defer db.Close()
	db.Find(&sources)

	return sources
}

func GetSubscribedNormalSources() []Source {
	var subscribedSources []Source
	sources := GetSources()
	for _, source := range sources {
		if source.IsSubscribed() && source.ErrorCount < config.ErrorThreshold {
			subscribedSources = append(subscribedSources, source)
		}
	}
	return subscribedSources
}

func (s *Source) IsSubscribed() bool {
	db := getConnect()
	defer db.Close()
	var sub Subscribe
	db.Where("source_id=?", s.ID).First(&sub)
	return sub.SourceID == s.ID
}

func (s *Source) GetNewContents() ([]Content, error) {
	var newContents []Content
	feed, err := rss.FetchByFunc(fetchFunc, s.Link)
	if err != nil {
		log.Println("Unable to make request: ", err, " ", s.Link)
		s.AddErrorCount()
		return nil, err
	}

	s.EraseErrorCount()

	items := feed.Items

	for _, item := range items {
		c, isBroad, _ := GenContentAndCheckByFeedItem(s, item)
		if !isBroad {
			newContents = append(newContents, *c)
		}
	}
	return newContents, nil
}

func GetSourcesByUserID(userID int64) ([]Source, error) {
	db := getConnect()
	defer db.Close()
	var sources []Source
	subs, err := GetSubsByUserID(userID)

	if err != nil {
		return nil, err
	}

	for _, sub := range subs {
		var source Source
		db.Where("id=?", sub.SourceID).First(&source)
		if source.ID == sub.SourceID {
			sources = append(sources, source)
		}
	}

	return sources, nil
}

func (s *Source) AddErrorCount() {
	s.ErrorCount++
	s.Save()
}

func (s *Source) EraseErrorCount() {
	s.ErrorCount = 0
	s.Save()
}

func (s *Source) Save() {
	db := getConnect()
	defer db.Close()
	db.Save(&s)
	return
}

func GetSourceById(id uint) (*Source, error) {
	db := getConnect()
	defer db.Close()
	var source Source

	if err := db.Where("id=?", id).First(&source); err.Error != nil {
		return nil, errors.New("未找到 RSS 源")
	}

	return &source, nil
}

func (s *Source) GetSubscribeNum() int {
	db := getConnect()
	defer db.Close()
	var subs []Subscribe
	db.Where("source_id=?", s.ID).Find(&subs)
	return len(subs)
}

func (s *Source) DeleteContents() {
	DeleteContentsBySourceID(s.ID)
}

func (s *Source) DeleteDueNoSubscriber() {
	s.DeleteContents()
	db := getConnect()
	defer db.Close()
	db.Delete(&s)
}
