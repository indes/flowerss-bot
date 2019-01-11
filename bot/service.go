package bot

import (
	"fmt"
	"github.com/indes/go-rssbot/bot/model"
	_ "github.com/indes/go-rssbot/bot/model"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

func registFeed(chat *telebot.Chat, url string) {
	msg, _ := B.Send(chat, "处理中...")
	log.Println(msg)

	k, _ := model.FindOrNewSourceByUrl(url)
	err := model.RegistFeed(chat.ID, k.ID)

	fmt.Println(k)

	if err == nil {
		msg, _ = B.Edit(msg, "订阅成功")
	} else {
		msg, _ = B.Edit(msg, "订阅失败")
	}
}
