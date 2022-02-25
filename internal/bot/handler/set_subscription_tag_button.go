package handler

import (
	"fmt"
	"strconv"
	"strings"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/model"
)

const (
	SetSubscriptionTagButtonUnique = "set_set_sub_tag_btn"
)

type SetSubscriptionTagButton struct {
	bot *tb.Bot
}

func NewSetSubscriptionTagButton(bot *tb.Bot) *SetSubscriptionTagButton {
	return &SetSubscriptionTagButton{bot: bot}
}

func (b *SetSubscriptionTagButton) CallbackUnique() string {
	return "\f" + SetSubscriptionTagButtonUnique
}

func (b *SetSubscriptionTagButton) Description() string {
	return ""
}

func (b *SetSubscriptionTagButton) feedSetAuth(c *tb.Callback) bool {
	data := strings.Split(c.Data, ":")
	subscriberID, _ := strconv.ParseInt(data[0], 10, 64)
	// 如果订阅者与按钮点击者id不一致，需要验证管理员权限
	if subscriberID != c.Sender.ID {
		channelChat, err := b.bot.ChatByID(subscriberID)
		if err != nil {
			return false
		}

		if !chat.IsChatAdmin(b.bot, channelChat, c.Sender.ID) {
			return false
		}
	}
	return true
}

func (b *SetSubscriptionTagButton) Handle(ctx tb.Context) error {
	c := ctx.Callback()
	// 权限验证
	if !b.feedSetAuth(c) {
		return ctx.Send("无权限")
	}
	data := strings.Split(c.Data, ":")
	ownID, _ := strconv.Atoi(data[0])
	sourceID, _ := strconv.Atoi(data[1])

	sub, err := model.GetSubscribeByUserIDAndSourceID(int64(ownID), uint(sourceID))
	if err != nil {
		return ctx.Send("系统错误，代码04")
	}
	msg := fmt.Sprintf(
		"请使用`/setfeedtag %d tags`命令为该订阅设置标签，tags为需要设置的标签，以空格分隔。（最多设置三个标签） \n"+
			"例如：`/setfeedtag %d 科技 苹果`",
		sub.ID, sub.ID,
	)
	return ctx.Edit(msg, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
}

func (b *SetSubscriptionTagButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
