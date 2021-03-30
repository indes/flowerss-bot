package main

import (
	"github.com/xos/rssbot/bot"
	"github.com/xos/rssbot/model"
	"github.com/xos/rssbot/task"
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
