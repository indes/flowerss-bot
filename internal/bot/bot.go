package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/handler"
	"github.com/indes/flowerss-bot/internal/bot/middleware"
	"github.com/indes/flowerss-bot/internal/bot/preview"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
)

type Bot struct {
	core *core.Core
	tb   *tb.Bot // telebot.Bot instance
}

func NewBot(core *core.Core) *Bot {
	log.Infof("init telegram bot, token %s, endpoint %s", config.BotToken, config.TelegramEndpoint)
	settings := tb.Settings{
		URL:    config.TelegramEndpoint,
		Token:  config.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: core.HttpClient().Client(),
	}

	logLevel := config.GetString("log.level")
	if strings.ToLower(logLevel) == "debug" {
		settings.Verbose = true
	}

	b := &Bot{
		core: core,
	}

	var err error
	b.tb, err = tb.NewBot(settings)
	if err != nil {
		log.Error(err)
		return nil
	}
	b.tb.Use(middleware.UserFilter(), middleware.PreLoadMentionChat(), middleware.IsChatAdmin())
	return b
}

func (b *Bot) registerCommands(appCore *core.Core) error {
	commandHandlers := []handler.CommandHandler{
		handler.NewStart(),
		handler.NewPing(b.tb),
		handler.NewAddSubscription(appCore),
		handler.NewRemoveSubscription(b.tb, appCore),
		handler.NewListSubscription(appCore),
		handler.NewRemoveAllSubscription(),
		handler.NewOnDocument(b.tb, appCore),
		handler.NewSet(b.tb, appCore),
		handler.NewSetFeedTag(appCore),
		handler.NewSetUpdateInterval(appCore),
		handler.NewExport(appCore),
		handler.NewImport(),
		handler.NewPauseAll(appCore),
		handler.NewActiveAll(appCore),
		handler.NewHelp(),
		handler.NewVersion(),
	}

	for _, h := range commandHandlers {
		b.tb.Handle(h.Command(), h.Handle, h.Middlewares()...)
	}

	ButtonHandlers := []handler.ButtonHandler{
		handler.NewRemoveAllSubscriptionButton(appCore),
		handler.NewCancelRemoveAllSubscriptionButton(),
		handler.NewSetFeedItemButton(b.tb, appCore),
		handler.NewRemoveSubscriptionItemButton(appCore),
		handler.NewNotificationSwitchButton(b.tb, appCore),
		handler.NewSetSubscriptionTagButton(b.tb),
		handler.NewTelegraphSwitchButton(b.tb, appCore),
		handler.NewSubscriptionSwitchButton(b.tb, appCore),
	}

	for _, h := range ButtonHandlers {
		b.tb.Handle(h, h.Handle, h.Middlewares()...)
	}

	var commands []tb.Command
	for _, h := range commandHandlers {
		if h.Description() == "" {
			continue
		}
		commands = append(commands, tb.Command{Text: h.Command(), Description: h.Description()})
	}
	log.Debugf("set bot command %+v", commands)
	if err := b.tb.SetCommands(commands); err != nil {
		return err
	}
	return nil
}

func (b *Bot) Run() error {
	if config.RunMode == config.TestMode {
		return nil
	}

	if err := b.registerCommands(b.core); err != nil {
		return err
	}
	log.Infof("bot start %s", config.AppVersionInfo())
	b.tb.Start()
	return nil
}

func (b *Bot) SourceUpdate(
	source *model.Source, newContents []*model.Content, subscribes []*model.Subscribe,
) {
	b.BroadcastNews(source, subscribes, newContents)
}

func (b *Bot) SourceUpdateError(source *model.Source) {
	b.BroadcastSourceError(source)
}

// BroadcastNews send new contents message to subscriber
func (b *Bot) BroadcastNews(source *model.Source, subs []*model.Subscribe, contents []*model.Content) {
	zap.S().Infow(
		"broadcast news",
		"fetcher id", source.ID,
		"fetcher title", source.Title,
		"subscriber count", len(subs),
		"new contents", len(contents),
	)

	for _, content := range contents {
		previewText := preview.TrimDescription(content.Description, config.PreviewText)

		for _, sub := range subs {
			tpldata := &config.TplData{
				SourceTitle:     source.Title,
				ContentTitle:    content.Title,
				RawLink:         content.RawLink,
				PreviewText:     previewText,
				TelegraphURL:    content.TelegraphURL,
				Tags:            sub.Tag,
				EnableTelegraph: sub.EnableTelegraph == 1 && content.TelegraphURL != "",
			}

			u := &tb.User{
				ID: sub.UserID,
			}
			o := &tb.SendOptions{
				DisableWebPagePreview: config.DisableWebPagePreview,
				ParseMode:             config.MessageMode,
				DisableNotification:   sub.EnableNotification != 1,
			}
			msg, err := tpldata.Render(config.MessageMode)
			if err != nil {
				zap.S().Errorw(
					"broadcast news error, tpldata.Render err",
					"error", err.Error(),
				)
				return
			}
			if _, err := b.tb.Send(u, msg, o); err != nil {

				if strings.Contains(err.Error(), "Forbidden") {
					zap.S().Errorw(
						"broadcast news error, bot stopped by user",
						"error", err.Error(),
						"user id", sub.UserID,
						"source id", sub.SourceID,
						"title", source.Title,
						"link", source.Link,
					)
					b.core.Unsubscribe(context.Background(), sub.UserID, sub.SourceID)
				}

				/*
					Telegram return error if markdown message has incomplete format.
					Print the msg to warn the user
					api error: Bad Request: can't parse entities: Can't find end of the entity starting at byte offset 894
				*/
				if strings.Contains(err.Error(), "parse entities") {
					zap.S().Errorw(
						"broadcast news error, markdown error",
						"markdown msg", msg,
						"error", err.Error(),
					)
				}
			}
		}
	}
}

// BroadcastSourceError send fetcher update error message to subscribers
func (b *Bot) BroadcastSourceError(source *model.Source) {
	subs, err := b.core.GetSourceAllSubscriptions(context.Background(), source.ID)
	if err != nil {
		log.Errorf("get subscriptions failed, %v", err)
	}
	var u tb.User
	for _, sub := range subs {
		message := fmt.Sprintf(
			"[%s](%s) 已经累计连续%d次更新失败，暂停更新", source.Title, source.Link, config.ErrorThreshold,
		)
		u.ID = sub.UserID
		_, _ = b.tb.Send(
			&u, message, &tb.SendOptions{
				ParseMode: tb.ModeMarkdown,
			},
		)
	}
}
