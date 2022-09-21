package handler

import (
	"context"
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/core"
)

type ActiveAll struct {
	core *core.Core
}

func NewActiveAll(core *core.Core) *ActiveAll {
	return &ActiveAll{core: core}
}

func (a *ActiveAll) Command() string {
	return "/activeall"
}

func (a *ActiveAll) Description() string {
	return "开启抓取订阅更新"
}

func (a *ActiveAll) Handle(ctx tb.Context) error {
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	subscribeUserID := ctx.Chat().ID
	if mentionChat != nil {
		subscribeUserID = mentionChat.ID
	}

	source, err := a.core.GetUserSubscribedSources(context.Background(), subscribeUserID)
	if err != nil {
		return ctx.Reply("系统错误")
	}

	for _, s := range source {
		err := a.core.EnableSourceUpdate(context.Background(), s.ID)
		if err != nil {
			return ctx.Reply("激活失败")
		}
	}

	reply := "订阅已全部开启"
	if mentionChat != nil {
		reply = fmt.Sprintf("频道 [%s](https://t.me/%s) 订阅已全部开启", mentionChat.Title, mentionChat.Username)
	}

	return ctx.Reply(
		reply, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (a *ActiveAll) Middlewares() []tb.MiddlewareFunc {
	return nil
}
