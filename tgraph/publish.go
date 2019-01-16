package tgraph

import (
	"fmt"
	"log"
)

func PublishHtml(sourceTitle string, title string, rawLink string, html string) (string, error) {
	//html = fmt.Sprintf(
	//	"<p>本文章由 <a href=\"https://github.com/indes/rssflow\">RSSFlow</a> 抓取自RSS，版权归<a href=\"\">源站点</a>所有。</p><hr>",
	//) + html + fmt.Sprintf(
	//	"<hr><p>本文章由 <a href=\"https://github.com/indes/rssflow\">RSSFlow</a> 抓取自RSS，版权归<a href=\"\">源站点</a>所有。</p><p>查看原文：<a href=\"%s\">%s - %s</p>",
	//	rawLink,
	//	title,
	//	sourceTitle,
	//)

	html = html + fmt.Sprintf(
		"<hr><p>本文章由 <a href=\"https://github.com/indes/rssflow\">RSSFlow</a> 抓取自RSS，版权归<a href=\"\">源站点</a>所有。</p><p>查看原文：<a href=\"%s\">%s - %s</p>",
		rawLink,
		title,
		sourceTitle,
	)
	if page, err := client.CreatePageWithHTML(title+" - "+sourceTitle, authorName, authorUrl, html, true); err == nil {
		log.Printf("Created telegraph page url: %s", page.URL)
		return page.URL, err
	} else {
		log.Printf("Create telegraph page error: %s", err)
		return "", nil
	}
}
