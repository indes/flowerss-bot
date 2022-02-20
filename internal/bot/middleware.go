package bot

import (
	"strconv"
	"strings"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/chat"
)

//feedSetAuth 验证订阅设置按钮点击者权限
func feedSetAuth(c *tb.Callback) bool {
	if !chat.IsChatAdmin(B, c.Message.Chat, c.Sender.ID) {
		return false
	}

	data := strings.Split(c.Data, ":")
	subscriberID, _ := strconv.ParseInt(data[0], 10, 64)
	// 如果订阅者与按钮点击者id不一致，需要验证管理员权限
	if subscriberID != c.Sender.ID {
		channelChat, err := B.ChatByID(subscriberID)
		if err != nil {
			return false
		}

		if !chat.IsChatAdmin(B, channelChat, c.Sender.ID) {
			return false
		}
	}
	return true
}
