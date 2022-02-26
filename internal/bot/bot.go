package bot

import (
	"time"

	"github.com/indes/flowerss-bot/internal/bot/handler"
	"github.com/indes/flowerss-bot/internal/bot/middleware"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/util"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
)

var (
	// B bot
	B *tb.Bot
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
	// create bot
	var err error
	B, err = tb.NewBot(
		tb.Settings{
			URL:     config.TelegramEndpoint,
			Token:   config.BotToken,
			Poller:  &tb.LongPoller{Timeout: 10 * time.Second},
			Client:  util.HttpClient,
			Verbose: true,
		},
	)
	B.Use(
		middleware.UserFilter(),
		middleware.PreLoadMentionChat(),
		middleware.IsChatAdmin(),
	)
	if err != nil {
		zap.S().Fatal(err)
		return
	}
}

//Start bot
func Start() {
	if config.RunMode == config.TestMode {
		return
	}

	zap.S().Infof("bot start %s", config.AppVersionInfo())
	setCommands()
	B.Start()
}

func setCommands() {
	commandHandlers := []handler.CommandHandler{
		handler.NewStart(),
		handler.NewPing(B),
		handler.NewAddSubscription(),
		handler.NewRemoveSubscription(B),
		handler.NewListSubscription(),
		handler.NewRemoveAllSubscription(),
		handler.NewOnDocument(B),
		handler.NewSet(B),
		handler.NewSetFeedTag(),
		handler.NewSetUpdateInterval(),
		handler.NewExport(),
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
		handler.NewRemoveAllSubscriptionButton(),
		handler.NewCancelRemoveAllSubscriptionButton(),
		handler.NewSetFeedItemButton(B),
		handler.NewRemoveSubscriptionItemButton(),
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
