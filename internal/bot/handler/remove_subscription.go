package handler

import (
	"context"
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
)

type RemoveSubscription struct {
	bot  *tb.Bot
	core *core.Core
}

func NewRemoveSubscription(bot *tb.Bot, core *core.Core) *RemoveSubscription {
	return &RemoveSubscription{
		bot:  bot,
		core: core,
	}
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

	source, err := s.core.GetSourceByURL(context.Background(), sourceURL)
	if err != nil {
		return ctx.Reply("获取订阅信息错误")
	}

	log.Infof("%d for [%d]%s unsubscribe %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
	if err := s.core.Unsubscribe(context.Background(), channelChat.ID, source.ID); err != nil {
		log.Errorf(
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
		sources, err := s.core.GetUserSubscribedSources(context.Background(), ctx.Chat().ID)
		if err != nil {
			return ctx.Reply("获取订阅列表失败")
		}

		if len(sources) == 0 {
			return ctx.Reply("没有订阅")
		}

		var unsubFeedItemButtons [][]tb.InlineButton
		for _, source := range sources {
			attachData := &session.Attachment{
				UserId:   ctx.Chat().ID,
				SourceId: uint32(source.ID),
			}

			data := session.Marshal(attachData)
			unsubFeedItemButtons = append(
				unsubFeedItemButtons, []tb.InlineButton{
					{
						Unique: RemoveSubscriptionItemButtonUnique,
						Text:   fmt.Sprintf("[%d] %s", source.ID, source.Title),
						Data:   data,
					},
				},
			)
		}
		return ctx.Reply("请选择你要退订的源", &tb.ReplyMarkup{InlineKeyboard: unsubFeedItemButtons})
	}

	if !chat.IsChatAdmin(s.bot, ctx.Chat(), ctx.Sender().ID) {
		return ctx.Reply("非管理员无法执行此操作")
	}

	source, err := s.core.GetSourceByURL(context.Background(), sourceURL)
	if err != nil {
		return ctx.Reply("未订阅该RSS源")
	}

	log.Infof("%d unsubscribe [%d]%s %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
	if err := s.core.Unsubscribe(context.Background(), ctx.Chat().ID, source.ID); err != nil {
		log.Errorf(
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
	core *core.Core
}

func NewRemoveSubscriptionItemButton(core *core.Core) *RemoveSubscriptionItemButton {
	return &RemoveSubscriptionItemButton{core: core}
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

	attachData, err := session.UnmarshalAttachment(ctx.Callback().Data)
	if err != nil {
		return ctx.Edit("退订错误！")
	}

	userID := attachData.GetUserId()
	sourceID := uint(attachData.GetSourceId())
	source, err := r.core.GetSource(context.Background(), sourceID)
	if err != nil {
		return ctx.Edit("退订错误！")
	}

	if err := r.core.Unsubscribe(context.Background(), userID, sourceID); err != nil {
		log.Errorf("unsubscribe data %s failed, %v", ctx.Callback().Data, err)
		return ctx.Edit("退订错误！")
	}

	rtnMsg := fmt.Sprintf("[%d] <a href=\"%s\">%s</a> 退订成功", sourceID, source.Link, source.Title)
	return ctx.Edit(rtnMsg, &tb.SendOptions{ParseMode: tb.ModeHTML})
}

func (r *RemoveSubscriptionItemButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
