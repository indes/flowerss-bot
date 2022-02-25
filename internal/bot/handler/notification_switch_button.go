package handler

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/model"
)

const (
	NotificationSwitchButtonUnique = "set_toggle_notice_btn"
)

type NotificationSwitchButton struct {
	bot *tb.Bot
}

func NewNotificationSwitchButton(bot *tb.Bot) *NotificationSwitchButton {
	return &NotificationSwitchButton{bot: bot}
}

func (b *NotificationSwitchButton) CallbackUnique() string {
	return "\f" + NotificationSwitchButtonUnique
}

func (b *NotificationSwitchButton) Description() string {
	return ""
}

func (b *NotificationSwitchButton) Handle(ctx tb.Context) error {
	c := ctx.Callback()
	if c == nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}

	data := strings.Split(c.Data, ":")
	subscriberID, _ := strconv.ParseInt(data[0], 10, 64)
	if subscriberID != c.Sender.ID {
		// 如果订阅者与按钮点击者id不一致，需要验证管理员权限
		channelChat, err := b.bot.ChatByID(subscriberID)
		if err != nil {
			return ctx.Respond(&tb.CallbackResponse{Text: "error"})
		}
		if !chat.IsChatAdmin(b.bot, channelChat, c.Sender.ID) {
			return ctx.Respond(&tb.CallbackResponse{Text: "error"})
		}
	}

	msg := strings.Split(c.Message.Text, "\n")
	subID, err := strconv.Atoi(strings.Split(msg[1], " ")[1])
	if err != nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}
	sub, err := model.GetSubscribeByID(subID)
	if sub == nil || err != nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}

	source, _ := model.GetSourceById(sub.SourceID)
	t := template.New("setting template")
	_, _ = t.Parse(feedSettingTmpl)

	err = sub.ToggleNotification()
	if err != nil {
		return ctx.Respond(&tb.CallbackResponse{Text: "error"})
	}
	sub.Save()
	text := new(bytes.Buffer)
	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})
	_ = ctx.Respond(&tb.CallbackResponse{Text: "修改成功"})
	return ctx.Edit(
		text.String(),
		&tb.SendOptions{ParseMode: tb.ModeHTML},
		&tb.ReplyMarkup{InlineKeyboard: genFeedSetBtn(c, sub, source)},
	)
}

func (b *NotificationSwitchButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
