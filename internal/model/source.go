package model

import (
	"fmt"
	"sort"

	"github.com/SlyMarbo/rss"
	"github.com/jinzhu/gorm"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/fetch"
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

func FindOrNewSourceByUrl(url string) (*Source, error) {
	var source Source

	if err := db.Where("link=?", url).Find(&source).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			source.Link = url

			// parsing task
			feed, err := rss.FetchByFunc(fetch.FetchFunc(httpClient), url)

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

func GetSources() (sources []*Source) {
	db.Find(&sources)
	return sources
}

func GetSubscribedNormalSources() []*Source {
	var subscribedSources []*Source
	sources := GetSources()
	for _, source := range sources {
		if source.IsSubscribed() && source.ErrorCount < config.ErrorThreshold {
			subscribedSources = append(subscribedSources, source)
		}
	}
	sort.SliceStable(
		subscribedSources, func(i, j int) bool {
			return subscribedSources[i].ID < subscribedSources[j].ID
		},
	)
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
