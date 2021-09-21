package main

import (
	"github.com/indes/flowerss-bot/internal/bot"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/task"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	model.InitDB()
	task.StartTasks()
	go handleSignal()
	bot.Start()
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	task.StopTasks()
	model.Disconnect()
	os.Exit(0)
}
