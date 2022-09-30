package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/indes/flowerss-bot/internal/bot"

	//"github.com/indes/flowerss-bot/internal/bot"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
)

func main() {
	model.InitDB()
	go handleSignal()

	appCore := core.NewCoreFormConfig()
	if err := appCore.Init(); err != nil {
		log.Fatal(err)
	}

	b := bot.NewBot(appCore)
	appCore.RegisterRssUpdateObserver(b)
	appCore.Run()
	b.Run()
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	model.Disconnect()
	os.Exit(0)
}
