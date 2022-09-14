package handler

import (
	"context"
	"strings"

	"github.com/spf13/cast"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/core"
)

type SetFeedTag struct {
	core *core.Core
}

func NewSetFeedTag(core *core.Core) *SetFeedTag {
	return &SetFeedTag{core: core}
}

func (s *SetFeedTag) Command() string {
	return "/setfeedtag"
}

func (s *SetFeedTag) Description() string {
	return "设置rss订阅标签"
}

func (s *SetFeedTag) getMessageWithoutMention(ctx tb.Context) string {
	mention := message.MentionFromMessage(ctx.Message())
	if mention == "" {
		return ctx.Message().Payload
	}
	return strings.Replace(ctx.Message().Payload, mention, "", -1)
}

func (s *SetFeedTag) Handle(ctx tb.Context) error {
	msg := s.getMessageWithoutMention(ctx)
	args := strings.Split(strings.TrimSpace(msg), " ")
	if len(args) < 1 {
		return ctx.Reply("/setfeedtag [sourceID] [tag1] [tag2] 设置订阅标签（最多设置三个Tag，以空格分割）")
	}

	// 截短参数
	if len(args) > 4 {
		args = args[:4]
	}

	sourceID := cast.ToUint(args[0])
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	subscribeUserID := ctx.Chat().ID
	if mentionChat != nil {
		subscribeUserID = mentionChat.ID
	}

	if err := s.core.SetSubscriptionTag(context.Background(), subscribeUserID, sourceID, args[1:]); err != nil {
		return ctx.Reply("订阅标签设置失败!")
	}
	return ctx.Reply("订阅标签设置成功!")
}

func (s *SetFeedTag) Middlewares() []tb.MiddlewareFunc {
	return nil
}
