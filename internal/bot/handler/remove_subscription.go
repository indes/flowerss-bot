package handler

import (
	"fmt"
	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/model"
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
)

type RemoveSubscription struct {
	bot *tb.Bot
}

func NewRemoveSubscription(bot *tb.Bot) *RemoveSubscription {
	return &RemoveSubscription{bot: bot}
}

func (r *RemoveSubscription) Command() string {
	return "/unsub"
}

func (r *RemoveSubscription) Description() string {
	return "退订RSS源"
}

func (r *RemoveSubscription) removeForChannel(ctx tb.Context, channelName string) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		return ctx.Send("频道退订请使用' /unsub @ChannelID URL ' 命令")
	}

	channelChat, err := r.bot.ChatByUsername(channelName)
	if err != nil {
		return ctx.Reply("获取频道信息错误")
	}

	if !chat.IsChatAdmin(r.bot, channelChat, ctx.Sender().ID) {
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
		fmt.Sprintf("频道 [%s](https://t.me/%s) 退订 [%s](%s) 成功",
			channelChat.Title, channelChat.Username, source.Title, source.Link),
		&tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown},
	)
}

func (r *RemoveSubscription) removeForChat(ctx tb.Context) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		return ctx.Reply("退订请使用' /unsub URL' 命令")
	}

	if !chat.IsChatAdmin(r.bot, ctx.Chat(), ctx.Sender().ID) {
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
	//Unsub by button
	//subs, err := model.GetSubsByUserID(m.Chat.ID)
	//if err != nil {
	//	return
	//}
	//
	//if len(subs) > 0 {
	//	unsubFeedItemBtns := [][]tb.InlineButton{}
	//
	//	for _, sub := range subs {
	//
	//		source, err := model.GetSourceById(sub.SourceID)
	//		if err != nil {
	//			return
	//		}
	//
	//		unsubFeedItemBtns = append(unsubFeedItemBtns, []tb.InlineButton{
	//			tb.InlineButton{
	//				Unique: "unsub_feed_item_btn",
	//				Text:   fmt.Sprintf("[%d] %s", sub.SourceID, source.Title),
	//				Data:   fmt.Sprintf("%d:%d:%d", sub.UserID, sub.ID, source.ID),
	//			},
	//		})
	//	}
	//
	//	_, _ = B.Send(m.Chat, "请选择你要退订的源", &tb.ReplyMarkup{
	//		InlineKeyboard: unsubFeedItemBtns,
	//	})
	//} else {
	//	_, _ = B.Send(m.Chat, "当前没有订阅源")
	//}
}

func (r *RemoveSubscription) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	if mention != "" {
		return r.removeForChannel(ctx, mention)
	}
	return r.removeForChat(ctx)
}

func (l *RemoveSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}
