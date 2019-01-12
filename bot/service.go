package bot

import (
	"fmt"
	"github.com/indes/go-rssbot/model"
	"gopkg.in/tucnak/telebot.v2"
)

func registFeed(chat *telebot.Chat, url string) {
	msg, _ := B.Send(chat, "处理中...")

	k, _ := model.FindOrNewSourceByUrl(url)
	err := model.RegistFeed(chat.ID, k.ID)

	fmt.Println(k)

	if err == nil {
		msg, _ = B.Edit(msg, fmt.Sprintf("<%s> 订阅成功", k.Title))
	} else {
		msg, _ = B.Edit(msg, "订阅失败")
	}
}

func SendError(c *telebot.Chat) {
	B.Send(c, "请输入正确的指令！")
}