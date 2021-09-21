package bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

//feedSetAuth 验证订阅设置按钮点击者权限
func feedSetAuth(c *tb.Callback) bool {
	if (c.Message.Chat.Type == tb.ChatGroup || c.Message.Chat.Type == tb.ChatSuperGroup) &&
		!userIsAdminOfGroup(c.Sender.ID, c.Message.Chat) {
		// check admin
		return false
	}

	data := strings.Split(c.Data, ":")
	subscriberID, _ := strconv.Atoi(data[0])
	// 如果订阅者与按钮点击者id不一致，需要验证管理员权限
	if subscriberID != c.Sender.ID {
		channelChat, err := B.ChatByID(fmt.Sprintf("%d", subscriberID))

		if err != nil {
			return false
		}

		if !UserIsAdminChannel(c.Sender.ID, channelChat) {
			return false
		}
	}

	return true
}

func checkPermit(userID int64, chatID int64) bool {
	// 个人用户
	if userID == chatID {
		return true
	}

	// 群组或频道
	chat, err := B.ChatByID(fmt.Sprintf("%d", chatID))

	if err != nil {
		return false
	}

	return checkPermitOfChat(userID, chat)
}

func checkPermitOfChat(userID int64, chat *tb.Chat) bool {
	if (chat.Type == tb.ChatGroup || chat.Type == tb.ChatSuperGroup) &&
		!userIsAdminOfGroup(int(userID), chat) {
		// check admin
		return false
	}
	return true
}
