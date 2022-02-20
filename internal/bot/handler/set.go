package handler

import (
	"fmt"

	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/model"

	tb "gopkg.in/telebot.v3"
)

type Set struct {
	bot *tb.Bot
}

func NewSet(bot *tb.Bot) *Set {
	return &Set{bot: bot}
}

func (s Set) Command() string {
	return "/set"
}

func (s Set) Description() string {
	return "设置订阅"
}

func (s Set) Handle(ctx tb.Context) error {
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	ownerID := ctx.Message().Chat.ID
	if mentionChat != nil {
		ownerID = mentionChat.ID
	}

	sources, err := model.GetSourcesByUserID(ownerID)
	if err != nil {
		return ctx.Reply("获取订阅失败")
	}
	if len(sources) <= 0 {
		return ctx.Reply("当前没有订阅")
	}

	// 配置按钮
	var replyButton []tb.ReplyButton
	replyKeys := [][]tb.ReplyButton{}
	setFeedItemBtns := [][]tb.InlineButton{}
	for _, source := range sources {
		// 添加按钮
		text := fmt.Sprintf("%s %s", source.Title, source.Link)
		replyButton = []tb.ReplyButton{
			tb.ReplyButton{Text: text},
		}
		replyKeys = append(replyKeys, replyButton)

		setFeedItemBtns = append(setFeedItemBtns, []tb.InlineButton{
			tb.InlineButton{
				Unique: "set_feed_item_btn",
				Text:   fmt.Sprintf("[%d] %s", source.ID, source.Title),
				Data:   fmt.Sprintf("%d:%d", ownerID, source.ID),
			},
		})
	}

	return ctx.Reply("请选择你要设置的源", &tb.ReplyMarkup{
		InlineKeyboard: setFeedItemBtns,
	})
}

func (s Set) Middlewares() []tb.MiddlewareFunc {
	return nil
}
