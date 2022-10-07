package tgraph

import (
	"html"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/indes/flowerss-bot/internal/log"
)

func PublishHtml(sourceTitle string, title string, rawLink string, htmlContent string) (string, error) {
	htmlContent = html.UnescapeString(htmlContent)
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	client := clientPool[rand.Intn(len(clientPool))]
	if page, err := client.CreatePageWithHTML(
		title+" - "+sourceTitle, sourceTitle, rawLink, htmlContent, true,
	); err == nil {
		zap.S().Infof("Created telegraph page url: %s", page.URL)
		return page.URL, err
	} else {
		log.Warnf("Create telegraph page failed, error: %s", err)
		return "", nil
	}
}
