package middleware

import (
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/config"
)

func UserFilter() tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			if len(config.AllowUsers) == 0 {
				return next(c)
			}
			userID := c.Sender().ID
			for _, allowUserID := range config.AllowUsers {
				if allowUserID == userID {
					return next(c)
				}
			}
			return fmt.Errorf("deny user %d", userID)
		}
	}
}
