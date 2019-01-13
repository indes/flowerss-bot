package main

import (
	"github.com/indes/rssflow/bot"
	"github.com/indes/rssflow/rss"
)

func main() {
	go rss.Update()
	bot.Start()
}
