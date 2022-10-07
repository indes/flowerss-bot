package handler

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
)

type AddSubscription struct {
	core *core.Core
}

func NewAddSubscription(core *core.Core) *AddSubscription {
	return &AddSubscription{
		core: core,
	}
}

func (a *AddSubscription) Command() string {
	return "/sub"
}

func (a *AddSubscription) Description() string {
	return "订阅RSS源"
}

func (a *AddSubscription) addSubscriptionForChat(ctx tb.Context) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		// 未附带链接，使用
		hint := fmt.Sprintf("请在命令后带上需要订阅的RSS URL，例如：%s https://justinpot.com/feed/", a.Command())
		return ctx.Send(hint, &tb.SendOptions{ReplyTo: ctx.Message()})
	}

	source, err := a.core.CreateSource(context.Background(), sourceURL)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("%s，订阅失败", err))
	}

	log.Infof("%d subscribe [%d]%s %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
	if err := a.core.AddSubscription(context.Background(), ctx.Chat().ID, source.ID); err != nil {
		if err == core.ErrSubscriptionExist {
			return ctx.Reply("已订阅该源，请勿重复订阅")
		}
		log.Errorf("add subscription user %d source %d failed %v", ctx.Chat().ID, source.ID, err)
		return ctx.Reply("订阅失败")
	}

	return ctx.Reply(
		fmt.Sprintf("[[%d]][%s](%s) 订阅成功", source.ID, source.Title, source.Link),
		&tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (a *AddSubscription) hasChannelPrivilege(bot *tb.Bot, channelChat *tb.Chat, opUserID int64, botID int64) (
	bool, error,
) {
	adminList, err := bot.AdminsOf(channelChat)
	if err != nil {
		zap.S().Error(err)
		return false, errors.New("获取频道信息失败")
	}

	senderIsAdmin := false
	botIsAdmin := false
	for _, admin := range adminList {
		if opUserID == admin.User.ID {
			senderIsAdmin = true
		}
		if botID == admin.User.ID {
			botIsAdmin = true
		}
	}

	return botIsAdmin && senderIsAdmin, nil
}

func (a *AddSubscription) addSubscriptionForChannel(ctx tb.Context, channelName string) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		return ctx.Send("频道订阅请使用' /sub @ChannelID URL ' 命令")
	}

	bot := ctx.Bot()
	channelChat, err := bot.ChatByUsername(channelName)
	if err != nil {
		return ctx.Reply("获取频道信息失败")
	}
	if channelChat.Type != tb.ChatChannel {
		return ctx.Reply("您或Bot不是频道管理员，无法设置订阅")
	}

	hasPrivilege, err := a.hasChannelPrivilege(bot, channelChat, ctx.Sender().ID, bot.Me.ID)
	if err != nil {
		return ctx.Reply(err.Error())
	}
	if !hasPrivilege {
		return ctx.Reply("您或Bot不是频道管理员，无法设置订阅")
	}

	source, err := a.core.CreateSource(context.Background(), sourceURL)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("%s，订阅失败", err))
	}

	log.Infof("%d subscribe [%d]%s %s", channelChat.ID, source.ID, source.Title, source.Link)
	if err := a.core.AddSubscription(context.Background(), channelChat.ID, source.ID); err != nil {
		if err == core.ErrSubscriptionExist {
			return ctx.Reply("已订阅该源，请勿重复订阅")
		}
		log.Errorf("add subscription user %d source %d failed %v", channelChat.ID, source.ID, err)
		return ctx.Reply("订阅失败")
	}

	return ctx.Reply(
		fmt.Sprintf("[[%d]] [%s](%s) 订阅成功", source.ID, source.Title, source.Link),
		&tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (a *AddSubscription) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	if mention != "" {
		// has mention, add subscription for channel
		return a.addSubscriptionForChannel(ctx, mention)
	}
	return a.addSubscriptionForChat(ctx)
}

func (a *AddSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}
