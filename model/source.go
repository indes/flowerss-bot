package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"unicode"

	"github.com/SlyMarbo/rss"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/log"
	"github.com/indes/flowerss-bot/util"
	"github.com/jinzhu/gorm"
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
	if err := db.Where("link=?", url).Find(&source).Error; err != nil {
		return nil, err
	}
	return &source, nil
}

func fetchFunc(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
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
	sort.SliceStable(subscribedSources, func(i, j int) bool {
		return subscribedSources[i].ID < subscribedSources[j].ID
	})
	return subscribedSources
}

func (s *Source) IsSubscribed() bool {
	var sub Subscribe
	db.Where("source_id=?", s.ID).First(&sub)
	return sub.SourceID == s.ID
}

func (s *Source) NeedUpdate() bool {
	var sub Subscribe
	db.Where("source_id=?", s.ID).First(&sub)
	sub.WaitTime += config.UpdateInterval
	if sub.Interval <= sub.WaitTime {
		sub.WaitTime = 0
		db.Save(&sub)
		return true
	} else {
		db.Save(&sub)
		return false
	}
}

func (s *Source) GetNewContents() ([]Content, error) {
	log.Debugw("fetch source updates",
		"source", s,
	)
	var newContents []Content
	feed, err := rss.FetchByFunc(fetchFunc, s.Link)
	if err != nil {
		log.Errorw("unable to fetch update", "error", err, "source", s)
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

	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].ID < sources[j].ID
	})

	return sources, nil
}

func GetErrorSourcesByUserID(userID int64) ([]Source, error) {
	var sources []Source
	subs, err := GetSubsByUserID(userID)

	if err != nil {
		return nil, err
	}

	for _, sub := range subs {
		var source Source
		db.Where("id=?", sub.SourceID).First(&source)
		if source.ID == sub.SourceID && source.ErrorCount >= config.ErrorThreshold {
			sources = append(sources, source)
		}
	}

	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].ID < sources[j].ID
	})

	return sources, nil
}

func ActiveSourcesByUserID(userID int64) error {
	subs, err := GetSubsByUserID(userID)

	if err != nil {
		return err
	}

	for _, sub := range subs {
		var source Source
		db.Where("id=?", sub.SourceID).First(&source)
		if source.ID == sub.SourceID {
			source.ErrorCount = 0
			db.Save(&source)
		}
	}

	return nil
}

func PauseSourcesByUserID(userID int64) error {
	subs, err := GetSubsByUserID(userID)

	if err != nil {
		return err
	}

	for _, sub := range subs {
		var source Source
		db.Where("id=?", sub.SourceID).First(&source)
		if source.ID == sub.SourceID {
			source.ErrorCount = config.ErrorThreshold + 1
			db.Save(&source)
		}
	}

	return nil
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
	db.Save(&s)
}

func GetSourceById(id uint) (*Source, error) {
	var source Source

	if err := db.Where("id=?", id).First(&source); err.Error != nil {
		return nil, errors.New("未找到 RSS 源")
	}

	return &source, nil
}

func (s *Source) GetSubscribeNum() int {
	var subs []Subscribe
	db.Where("source_id=?", s.ID).Find(&subs)
	return len(subs)
}

func (s *Source) DeleteContents() {
	DeleteContentsBySourceID(s.ID)
}

func (s *Source) DeleteDueNoSubscriber() {
	s.DeleteContents()
	db.Delete(&s)
}
