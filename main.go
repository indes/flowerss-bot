package main

import (
	"github.com/indes/rssflow/bot"
)

func main() {
	go task.Update()
	bot.Start()
}
