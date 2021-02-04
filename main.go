package main

import (
	"github.com/yangon99/flowerss-bot/bot"
	"github.com/yangon99/flowerss-bot/model"
	"github.com/yangon99/flowerss-bot/task"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	model.InitDB()
	go task.Update()
	go handleSignal()
	bot.Start()
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	model.Disconnect()
	os.Exit(0)
}
