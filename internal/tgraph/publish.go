package tgraph

import (
	"fmt"
	"html"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

func PublishHtml(sourceTitle string, title string, rawLink string, htmlContent string) (string, error) {
	//html = fmt.Sprintf(
	//	"<p>本文章由 <a href=\"https://t.me/freeinformation_bot\">無料案內所</a> 抓取，版權歸<a href=\"\">原網站</a>所有。</p><hr>",
	//) + html + fmt.Sprintf(
	//	"<hr><p>本文章由 <a href=\"https://t.me/freeinformation_bot\">無料案內所</a> 抓取，版權歸<a href=\"\">原網站</a>所有。</p><p>查看原文：<a href=\"%s\">%s - %s</p>",
	//	rawLink,
	//	title,
	//	sourceTitle,
	//)

	htmlContent = html.UnescapeString(htmlContent) + fmt.Sprintf(
		"<hr><p>本文章由 <a href=\"https://t.me/freeinformation_bot\">無料案內所</a> 抓取，版權歸<a href=\"\">原網站</a>所有。</p><p>查看原文：<a href=\"%s\">%s - %s</p>",
		rawLink,
		title,
		sourceTitle,
	)
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	client := clientPool[rand.Intn(len(clientPool))]

	if page, err := client.CreatePageWithHTML(title+" - "+sourceTitle, sourceTitle, rawLink, htmlContent, true); err == nil {
		zap.S().Infof("Created telegraph page url: %s", page.URL)
		return page.URL, err
	} else {
		zap.S().Warnf("Create telegraph page failed, error: %s", err)
		return "", nil
	}
}
