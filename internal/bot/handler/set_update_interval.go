package handler

import (
	"context"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
)

type SetUpdateInterval struct {
	core *core.Core
}

func NewSetUpdateInterval(core *core.Core) *SetUpdateInterval {
	return &SetUpdateInterval{core: core}
}

func (s *SetUpdateInterval) Command() string {
	return "/setinterval"
}

func (s *SetUpdateInterval) Description() string {
	return "设置订阅刷新频率"
}

func (s *SetUpdateInterval) getMessageWithoutMention(ctx tb.Context) string {
	mention := message.MentionFromMessage(ctx.Message())
	if mention == "" {
		return ctx.Message().Payload
	}
	return strings.Replace(ctx.Message().Payload, mention, "", -1)
}

func (s *SetUpdateInterval) Handle(ctx tb.Context) error {
	msg := s.getMessageWithoutMention(ctx)
	args := strings.Split(strings.TrimSpace(msg), " ")
	if len(args) < 2 {
		return ctx.Reply("/setinterval [interval] [sourceID] 设置订阅刷新频率（可设置多个sub id，以空格分割）")
	}

	interval, err := strconv.Atoi(args[0])
	if interval <= 0 || err != nil {
		return ctx.Reply("请输入正确的抓取频率")
	}

	subscribeUserID := ctx.Message().Chat.ID
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	if mentionChat != nil {
		subscribeUserID = mentionChat.ID
	}

	for _, id := range args[1:] {
		sourceID := cast.ToUint(id)
		if err := s.core.SetSubscriptionInterval(
			context.Background(), subscribeUserID, sourceID, interval,
		); err != nil {
			log.Errorf("SetSubscriptionInterval failed, %v", err)
			return ctx.Reply("抓取频率设置失败!")
		}
	}
	return ctx.Reply("抓取频率设置成功!")
}

func (s *SetUpdateInterval) Middlewares() []tb.MiddlewareFunc {
	return nil
}
