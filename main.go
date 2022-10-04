package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/indes/flowerss-bot/internal/bot"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
)

func main() {
	model.InitDB()

	appCore := core.NewCoreFormConfig()
	if err := appCore.Init(); err != nil {
		log.Fatal(err)
	}
	go handleSignal(appCore)
	b := bot.NewBot(appCore)
	appCore.RegisterRssUpdateObserver(b)
	appCore.Run()
	b.Run()
}

func handleSignal(appCore *core.Core) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	appCore.Stop()
	model.Disconnect()
	os.Exit(0)
}
