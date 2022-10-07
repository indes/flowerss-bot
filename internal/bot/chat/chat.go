package chat

import (
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/log"
)

func IsChatAdmin(bot *tb.Bot, chat *tb.Chat, userID int64) bool {
	if chat == nil || bot == nil {
		log.Errorf("chat or bot is nil, chat %v bot %v", chat, bot)
		return false
	}

	if chat.Type == tb.ChatPrivate {
		return true
	}

	admins, err := bot.AdminsOf(chat)
	if err != nil {
		log.Warnf("get admins of chat %v failed, %v", chat.ID, err)
		return false
	}

	for _, admin := range admins {
		if userID != admin.User.ID {
			continue
		}
		return true
	}
	return false
}
