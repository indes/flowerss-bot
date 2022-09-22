package model

import (
	"errors"
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

func GetSubscriberBySource(s *Source) []*Subscribe {
	if s == nil {
		return []*Subscribe{}
	}

	var subs []*Subscribe

	db.Where("source_id=?", s.ID).Find(&subs)
	return subs
}

func (s *Subscribe) Unsub() error {
	if s.ID == 0 {
		return errors.New("can't delete 0 subscribe")
	}

	return db.Delete(&s).Error
}
