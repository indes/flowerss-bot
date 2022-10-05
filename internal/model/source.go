package model

import (
	"sort"

	"github.com/indes/flowerss-bot/internal/config"
)

type Source struct {
	ID         uint `gorm:"primary_key;AUTO_INCREMENT"`
	Link       string
	Title      string
	ErrorCount uint
	Content    []Content
	EditTime
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
