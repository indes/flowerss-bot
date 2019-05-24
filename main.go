package main

import (
	"github.com/indes/flowerss-bot/bot"
	"github.com/indes/flowerss-bot/task"
)

func main() {
	go task.Update()
	bot.Start()
}
