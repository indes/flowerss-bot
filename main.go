package main

import (
	"github.com/indes/rssflow/bot"
	"github.com/indes/rssflow/task"
)

func main() {
	go task.Update()
	bot.Start()
}
