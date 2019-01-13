package model

import (
	"github.com/SlyMarbo/rss"
	"github.com/indes/rssflow/tgraph"
	"strings"
)

type Content struct {
	SourceID     uint
	HashID       string `gorm:"primary_key"`
	RawID        string
	RawLink      string
	Title        string
	TelegraphUrl string
	EditTime
}

func getContentByFeedItem(sid uint, item *rss.Item) (Content, error) {
	tgpUrl := ""

	html := item.Content
	if html == "" {
		html = item.Summary
	}

	html = strings.Replace(html, "<![CDATA[", "", -1)
	html = strings.Replace(html, "]]>", "", -1)
	tgpUrl = tgraph.PublishItem(item.Title, html)

	var c = Content{
		Title:        item.Title,
		SourceID:     sid,
		RawID:        item.ID,
		HashID:       genHashID(sid, item.ID),
		TelegraphUrl: tgpUrl,
		RawLink:      item.Link,
	}

	return c, nil
}
