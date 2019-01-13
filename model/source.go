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
		c, _ := getContentByFeedItem(s, item)
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

func GetSources() []Source {
	var sources []Source

	db := getConnect()
	defer db.Close()
	db.Find(&sources)
	return sources
}

func (s *Source) GetNewContents() ([]Content, error) {
	var contents []Content
	feed, err := rss.Fetch(s.Link)
	if err != nil {
		log.Println("Unable to make request: ", err)
		return nil, err
	}

	items := feed.Items

	for _, item := range items {
		c, isBroad, _ := GenContentAndCheckByFeedItem(s, item)
		if isBroad {
			subs := getSubscriberBySource(s)
			log.Println(subs)
		}
		log.Println(c, isBroad)
	}
	return contents, nil
}
