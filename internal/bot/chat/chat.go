package chat

import (
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
)

func IsChatAdmin(bot *tb.Bot, chat *tb.Chat, userID int64) bool {
	if chat == nil || bot == nil {
		zap.S().Errorf("chat or bot is nil, chat %v bot %v", chat, bot)
		return false
	}

	if chat.Type == tb.ChatPrivate {
		return true
	}

	admins, err := bot.AdminsOf(chat)
	if err != nil {
		zap.S().Warnf("get admins of chat %v failed, %v", chat.ID, err)
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
