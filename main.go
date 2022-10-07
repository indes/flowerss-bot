package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/indes/flowerss-bot/internal/bot"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/scheduler"
)

func main() {
	appCore := core.NewCoreFormConfig()
	if err := appCore.Init(); err != nil {
		log.Fatal(err)
	}
	go handleSignal()
	b := bot.NewBot(appCore)

	task := scheduler.NewRssTask(appCore)
	task.Register(b)
	task.Start()
	b.Run()
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	os.Exit(0)
}
