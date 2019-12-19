package model

import (
	"github.com/SlyMarbo/rss"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/tgraph"
	"strings"
)

type Content struct {
	SourceID     uint
	HashID       string `gorm:"primary_key"`
	RawID        string
	RawLink      string
	Title        string
	Description  string `gorm:"-"` //ignore to db
	TelegraphUrl string
	EditTime
}

func getContentByFeedItem(source *Source, item *rss.Item) (Content, error) {
	tgpUrl := ""

	html := item.Content
	if html == "" {
		html = item.Summary
	}

	html = strings.Replace(html, "<![CDATA[", "", -1)
	html = strings.Replace(html, "]]>", "", -1)

	if config.EnableTelegraph && len([]rune(html)) > config.PreviewText {
		tgpUrl = PublishItem(source, item, html)
	}

	var c = Content{
		Title:        strings.Trim(item.Title, " "),
		Description:  html, //replace all kinds of <br> tag
		SourceID:     source.ID,
		RawID:        item.ID,
		HashID:       genHashID(source.Link, item.ID),
		TelegraphUrl: tgpUrl,
		RawLink:      item.Link,
	}

	return c, nil
}

func GenContentAndCheckByFeedItem(s *Source, item *rss.Item) (*Content, bool, error) {
	db := getConnect()
	defer db.Close()
	var (
		content   Content
		isBroaded bool
	)

	hashID := genHashID(s.Link, item.ID)
	db.Where("hash_id=?", hashID).First(&content)
	if content.HashID == "" {
		isBroaded = false
		content, _ = getContentByFeedItem(s, item)
		db.Save(&content)
	} else {
		isBroaded = true
	}

	return &content, isBroaded, nil
}

func DeleteContentsBySourceID(sid uint) {
	db := getConnect()
	defer db.Close()
	db.Delete(Content{}, "source_id = ?", sid)
}

func PublishItem(source *Source, item *rss.Item, html string) string {
	url, _ := tgraph.PublishHtml(source.Title, item.Title, item.Link, html)
	return url
}
