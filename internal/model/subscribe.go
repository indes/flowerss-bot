package model

import (
	"errors"
	"strings"

	"github.com/indes/flowerss-bot/internal/config"
)

const (
	MaxSubscribeTagLength = 250
)

type Subscribe struct {
	ID                 uint `gorm:"primary_key;AUTO_INCREMENT"`
	UserID             int64
	SourceID           uint
	EnableNotification int
	EnableTelegraph    int
	Tag                string
	Interval           int
	WaitTime           int
	EditTime
}

func GetSubscribeByUserIDAndSourceID(userID int64, sourceID uint) (*Subscribe, error) {
	var sub Subscribe
	db.Where("user_id=? and source_id=?", userID, sourceID).First(&sub)
	if sub.UserID != int64(userID) {
		return nil, errors.New("未订阅该RSS源")
	}
	return &sub, nil
}

func GetSubscriberBySource(s *Source) []*Subscribe {
	if s == nil {
		return []*Subscribe{}
	}

	var subs []*Subscribe

	db.Where("source_id=?", s.ID).Find(&subs)
	return subs
}

func GetSubsByUserID(userID int64) ([]Subscribe, error) {
	var subs []Subscribe
	db.Where("user_id=?", userID).Order("id").Find(&subs)
	return subs, nil
}

func GetSubscribeByID(id int) (*Subscribe, error) {
	var sub Subscribe
	err := db.Where("id=?  ", id).First(&sub).Error
	return &sub, err
}

func (s *Subscribe) ToggleNotification() error {
	if s.EnableNotification != 1 {
		s.EnableNotification = 1
	} else {
		s.EnableNotification = 0
	}
	return nil
}

func (s *Subscribe) ToggleTelegraph() error {
	if s.EnableTelegraph != 1 {
		s.EnableTelegraph = 1
	} else {
		s.EnableTelegraph = 0
	}
	return nil
}

func (s *Source) ToggleEnabled() error {
	if s.ErrorCount >= config.ErrorThreshold {
		s.ErrorCount = 0
	} else {
		s.ErrorCount = config.ErrorThreshold
	}

	///TODO a hack for save source changes
	s.Save()

	return nil
}

func (s *Subscribe) SetTag(tags []string) error {
	defer s.Save()

	tagStr := strings.Join(tags, " #")

	s.Tag = "#" + tagStr
	return nil
}

func (s *Subscribe) SetInterval(interval int) error {
	defer s.Save()
	s.Interval = interval
	return nil
}

func (s *Subscribe) Unsub() error {
	if s.ID == 0 {
		return errors.New("can't delete 0 subscribe")
	}

	return db.Delete(&s).Error
}

func (s *Subscribe) Save() {
	db.Save(&s)
}
