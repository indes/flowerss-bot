package bot

import (
	"fmt"
	"github.com/indes/rssflow/model"
	"gopkg.in/tucnak/telebot.v2"
	"log"
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

func BroadNews(subs []model.Subscribe, contents []model.Content) error {
	//bot.Start()
	var u telebot.User
	log.Println("Subs Len: ", len(subs), " Contents Len: ", len(contents))
	for _, content := range contents {
		for _, sub := range subs {
			u.ID = int(sub.UserID)
			B.Send(&u, content.TelegraphUrl)
		}
	}
	return nil
}
