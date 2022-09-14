package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/indes/flowerss-bot/internal/bot/handler"
	"github.com/indes/flowerss-bot/internal/bot/middleware"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/pkg/client"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
)

var (
	// B bot
	B *tb.Bot

	Core *core.Core
)

func init() {
	if config.RunMode == config.TestMode {
		return
	}
	zap.S().Infow(
		"init telegram bot",
		"token", config.BotToken,
		"endpoint", config.TelegramEndpoint,
	)

	clientOpts := []client.HttpClientOption{
		client.WithTimeout(10 * time.Second),
	}
	if config.Socks5 != "" {
		clientOpts = append(clientOpts, client.WithProxyURL(fmt.Sprintf("socks5://%s", config.Socks5)))
	}
	httpClient := client.NewHttpClient(clientOpts...)

	settings := tb.Settings{
		URL:    config.TelegramEndpoint,
		Token:  config.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: httpClient.Client(),
	}

	logLevel := config.GetString("log.level")
	if strings.ToLower(logLevel) == "debug" {
		settings.Verbose = true
	}

	var err error
	B, err = tb.NewBot(settings)
	if err != nil {
		zap.S().Fatal(err)
		return
	}
	B.Use(middleware.UserFilter(), middleware.PreLoadMentionChat(), middleware.IsChatAdmin())
}

// Start bot
func Start(appCore *core.Core) {
	if config.RunMode == config.TestMode {
		return
	}

	Core = appCore
	zap.S().Infof("bot start %s", config.AppVersionInfo())
	setCommands()
	B.Start()
}

func setCommands() {
	commandHandlers := []handler.CommandHandler{
		handler.NewStart(),
		handler.NewPing(B),
		handler.NewAddSubscription(Core),
		handler.NewRemoveSubscription(B, Core),
		handler.NewListSubscription(Core),
		handler.NewRemoveAllSubscription(),
		handler.NewOnDocument(B, Core),
		handler.NewSet(B, Core),
		handler.NewSetFeedTag(Core),
		handler.NewSetUpdateInterval(),
		handler.NewExport(Core),
		handler.NewImport(),
		handler.NewPauseAll(),
		handler.NewActiveAll(),
		handler.NewHelp(),
		handler.NewVersion(),
	}

	for _, h := range commandHandlers {
		B.Handle(h.Command(), h.Handle, h.Middlewares()...)
	}

	ButtonHandlers := []handler.ButtonHandler{
		handler.NewRemoveAllSubscriptionButton(Core),
		handler.NewCancelRemoveAllSubscriptionButton(),
		handler.NewSetFeedItemButton(B, Core),
		handler.NewRemoveSubscriptionItemButton(Core),
		handler.NewNotificationSwitchButton(B),
		handler.NewSetSubscriptionTagButton(B),
		handler.NewTelegraphSwitchButton(B),
		handler.NewSubscriptionSwitchButton(B),
	}

	for _, h := range ButtonHandlers {
		B.Handle(h, h.Handle, h.Middlewares()...)
	}

	var commands []tb.Command
	for _, h := range commandHandlers {
		if h.Description() == "" {
			continue
		}
		commands = append(commands, tb.Command{Text: h.Command(), Description: h.Description()})
	}
	zap.S().Debugf("set bot command %+v", commands)
	if err := B.SetCommands(commands); err != nil {
		zap.S().Errorw("set bot commands failed", "error", err.Error())
	}
}
