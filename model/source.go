package model

import (
	"github.com/SlyMarbo/rss"
	"log"
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
	for _, item := range items {
		c, _ := getContentByFeedItem(s.ID, item)
		s.Content = append(s.Content, c)
	}

	return nil
}

func FindOrNewSourceByUrl(url string) (*Source, error) {
	var source Source
	db := getConnect()
	defer db.Close()

	if err := db.Where("link=?", url).Find(&source).Error; err != nil {
		if err.Error() == "record not found" {
			source.Link = url

			// parsing rss
			feed, err := rss.Fetch(url)
			if err != nil {
				log.Println("Unable to make request: ", err)
				return nil, err
			}

			source.Title = feed.Title
			source.ErrorCount = 0

			// Get contents and insert
			items := feed.Items
			source.appendContents(items)
			db.Create(&source)
			return &source, nil
		}
		return nil, err
	}

	return &source, nil
}