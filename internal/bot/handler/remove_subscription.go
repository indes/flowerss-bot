package handler

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/model"
)

type RemoveSubscription struct {
	bot *tb.Bot
}

func NewRemoveSubscription(bot *tb.Bot) *RemoveSubscription {
	return &RemoveSubscription{bot: bot}
}

func (s *RemoveSubscription) Command() string {
	return "/unsub"
}

func (s *RemoveSubscription) Description() string {
	return "退订RSS源"
}

func (s *RemoveSubscription) removeForChannel(ctx tb.Context, channelName string) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		return ctx.Send("频道退订请使用' /unsub @ChannelID URL ' 命令")
	}

	channelChat, err := s.bot.ChatByUsername(channelName)
	if err != nil {
		return ctx.Reply("获取频道信息错误")
	}

	if !chat.IsChatAdmin(s.bot, channelChat, ctx.Sender().ID) {
		return ctx.Reply("非频道管理员无法执行此操作")
	}

	source, _ := model.GetSourceByUrl(sourceURL)
	sub, err := model.GetSubByUserIDAndURL(channelChat.ID, sourceURL)
	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Send(
				fmt.Sprintf("频道 [%s](https://t.me/%s) 未订阅该RSS源", channelChat.Title, channelChat.Username),
				&tb.SendOptions{
					DisableWebPagePreview: true,
					ParseMode:             tb.ModeMarkdown,
				},
			)
		}
		return ctx.Reply("退订失败")
	}
	zap.S().Infof("%d for [%d]%s unsubscribe %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
	if err := sub.Unsub(); err != nil {
		zap.S().Errorf(
			"%d for [%d]%s unsubscribe %s failed, %v",
			ctx.Chat().ID, source.ID, source.Title, source.Link, err,
		)
		return ctx.Reply("退订失败")
	}
	return ctx.Send(
		fmt.Sprintf(
			"频道 [%s](https://t.me/%s) 退订 [%s](%s) 成功",
			channelChat.Title, channelChat.Username, source.Title, source.Link,
		),
		&tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown},
	)
}

func (s *RemoveSubscription) removeForChat(ctx tb.Context) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		subs, err := model.GetSubsByUserID(ctx.Chat().ID)
		if err != nil {
			return ctx.Reply("获取订阅列表失败")
		}

		if len(subs) == 0 {
			return ctx.Reply("没有订阅")
		}

		var unsubFeedItemButtons [][]tb.InlineButton
		for _, sub := range subs {
			source, err := model.GetSourceById(sub.SourceID)
			if err != nil {
				return ctx.Reply("获取订阅列表失败")
			}
			unsubFeedItemButtons = append(
				unsubFeedItemButtons, []tb.InlineButton{
					{
						Unique: "unsub_feed_item_btn",
						Text:   fmt.Sprintf("[%d] %s", sub.SourceID, source.Title),
						Data:   fmt.Sprintf("%d:%d:%d", sub.UserID, sub.ID, source.ID),
					},
				},
			)
		}
		return ctx.Reply("请选择你要退订的源", &tb.ReplyMarkup{InlineKeyboard: unsubFeedItemButtons})
	}

	if !chat.IsChatAdmin(s.bot, ctx.Chat(), ctx.Sender().ID) {
		return ctx.Reply("非管理员无法执行此操作")
	}

	source, err := model.GetSourceByUrl(sourceURL)
	if err != nil || source == nil {
		return ctx.Reply("未订阅该RSS源")
	}

	zap.S().Infof("%d unsubscribe [%d]%s %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
	if err := model.UnsubByUserIDAndSource(ctx.Chat().ID, source); err != nil {
		zap.S().Errorf(
			"%d for [%d]%s unsubscribe %s failed, %v",
			ctx.Chat().ID, source.ID, source.Title, source.Link, err,
		)
		return ctx.Reply("退订失败")
	}
	return ctx.Send(
		fmt.Sprintf("[%s](%s) 退订成功！", source.Title, source.Link),
		&tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown},
	)
}

func (s *RemoveSubscription) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	if mention != "" {
		return s.removeForChannel(ctx, mention)
	}
	return s.removeForChat(ctx)
}

func (s *RemoveSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}

const (
	RemoveSubscriptionItemButtonUnique = "unsub_feed_item_btn"
)

type RemoveSubscriptionItemButton struct {
}

func NewRemoveSubscriptionItemButton() *RemoveSubscriptionItemButton {
	return &RemoveSubscriptionItemButton{}
}

func (r *RemoveSubscriptionItemButton) CallbackUnique() string {
	return "\f" + RemoveSubscriptionItemButtonUnique
}

func (r *RemoveSubscriptionItemButton) Description() string {
	return ""
}

func (r *RemoveSubscriptionItemButton) Handle(ctx tb.Context) error {
	if ctx.Callback() == nil {
		return ctx.Edit("内部错误！")
	}

	data := strings.Split(ctx.Callback().Data, ":")
	if len(data) != 3 {
		return ctx.Edit("退订错误！")
	}

	userID, _ := strconv.Atoi(data[0])
	subID, _ := strconv.Atoi(data[1])
	sourceID, _ := strconv.Atoi(data[2])
	source, err := model.GetSourceById(uint(sourceID))
	if err != nil {
		return ctx.Edit("退订错误！")
	}

	if err := model.UnsubByUserIDAndSubID(int64(userID), uint(subID)); err != nil {
		return ctx.Edit("退订错误！")
	}

	rtnMsg := fmt.Sprintf("[%d] <a href=\"%s\">%s</a> 退订成功", sourceID, source.Link, source.Title)
	return ctx.Edit(rtnMsg, &tb.SendOptions{ParseMode: tb.ModeHTML})
}

func (r *RemoveSubscriptionItemButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
