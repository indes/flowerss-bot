package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/indes/flowerss-bot/internal/bot"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/task"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	model.InitDB()
	task.StartTasks()
	go handleSignal()

	appCore := core.NewCoreFormConfig()
	if err := appCore.Init(); err != nil {
		log.Fatal(err)
	}

	bot.Start(appCore)
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	task.StopTasks()
	model.Disconnect()
	os.Exit(0)
}
