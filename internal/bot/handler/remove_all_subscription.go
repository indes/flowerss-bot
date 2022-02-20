package handler

import (
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/model"
)

type RemoveAllSubscription struct {
}

func NewRemoveAllSubscription() *RemoveAllSubscription {
	return &RemoveAllSubscription{}
}

func (r RemoveAllSubscription) Command() string {
	return "/unsuball"
}

func (r RemoveAllSubscription) Description() string {
	return "取消所有订阅"
}

func (r RemoveAllSubscription) Handle(ctx tb.Context) error {
	reply := "是否退订当前用户的所有订阅？"
	confirmKeys := [][]tb.InlineButton{}
	confirmKeys = append(
		confirmKeys, []tb.InlineButton{
			tb.InlineButton{
				Unique: UnSubAllButtonUnique,
				Text:   "确认",
			},
			tb.InlineButton{
				Unique: CancelUnSubAllButtonUnique,
				Text:   "取消",
			},
		},
	)
	return ctx.Reply(reply, &tb.ReplyMarkup{InlineKeyboard: confirmKeys})
}

func (r RemoveAllSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}

const (
	UnSubAllButtonUnique       = "unsub_all_confirm_btn"
	CancelUnSubAllButtonUnique = "unsub_all_cancel_btn"
)

type RemoveAllSubscriptionButton struct {
}

func (r *RemoveAllSubscriptionButton) CallbackUnique() string {
	return "\f" + UnSubAllButtonUnique
}

func (r *RemoveAllSubscriptionButton) Description() string {
	return ""
}

func (r *RemoveAllSubscriptionButton) Handle(ctx tb.Context) error {
	success, fail, err := model.UnsubAllByUserID(ctx.Sender().ID)
	if err != nil {
		return ctx.Edit("退订失败")
	}
	return ctx.Edit(fmt.Sprintf("退订成功：%d\n退订失败：%d", success, fail))
}

func (r *RemoveAllSubscriptionButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}

type CancelRemoveAllSubscriptionButton struct {
}

func (r *CancelRemoveAllSubscriptionButton) CallbackUnique() string {
	return "\f" + CancelUnSubAllButtonUnique
}

func (r *CancelRemoveAllSubscriptionButton) Description() string {
	return ""
}

func (r *CancelRemoveAllSubscriptionButton) Handle(ctx tb.Context) error {
	return ctx.Edit("操作取消")
}

func (r *CancelRemoveAllSubscriptionButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
