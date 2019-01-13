package model

import (
	"github.com/SlyMarbo/rss"
)

type Source struct {
	ID         uint `gorm:"primary_key";"AUTO_INCREMENT"`
	Link       string
	Title      string
	ErrorCount uint
	Content    []Content
	EditTime
}

func (s *Source) appendContents(items []*rss.Item) error {
	//htmlRegexp := regexp.MustCompile(`^$`)

	for _, item := range items {

		c, _ := getContentByFeedItem(s.ID, item)
		s.Content = append(s.Content, c)
	}

	return nil
}
