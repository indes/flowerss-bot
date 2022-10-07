package middleware

import (
	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/log"

	tb "gopkg.in/telebot.v3"
)

func PreLoadMentionChat() tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			mention := message.MentionFromMessage(c.Message())
			if mention != "" {
				chat, err := c.Bot().ChatByUsername(mention)
				if err != nil {
					log.Errorf("pre load mention %s chat failed, %v", mention, err)
					return next(c)
				}
				c.Set(session.StoreKeyMentionChat.String(), chat)
			}
			return next(c)
		}
	}
}
