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
	TelegraphSwitchButtonUnique = "set_toggle_telegraph_btn"
)

type TelegraphSwitchButton struct {
	bot *tb.Bot
}

func NewTelegraphSwitchButton(bot *tb.Bot) *TelegraphSwitchButton {
	return &TelegraphSwitchButton{bot: bot}
}

func (b *TelegraphSwitchButton) CallbackUnique() string {
	return "\f" + TelegraphSwitchButtonUnique
}

func (b *TelegraphSwitchButton) Description() string {
	return ""
}

func (b *TelegraphSwitchButton) Handle(ctx tb.Context) error {
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

	err = sub.ToggleTelegraph()
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

func (b *TelegraphSwitchButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
