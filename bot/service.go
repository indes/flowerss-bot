package bot

import (
	"fmt"
	"github.com/indes/rssflow/model"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

func registFeed(chat *telebot.Chat, url string) {
	msg, _ := B.Send(chat, "处理中...")

	source, err := model.FindOrNewSourceByUrl(url)

	if err != nil {
		msg, _ = B.Edit(msg, fmt.Sprintf("%s，订阅失败", err))
		return
	}

	err = model.RegistFeed(chat.ID, source.ID)
	log.Printf("%d subscribe [%d]%s %s", chat.ID, source.ID, source.Title, source.Link)

	if err == nil {
		_, _ = B.Edit(msg, fmt.Sprintf(" [%s](%s) 订阅成功", source.Title, source.Link),
			&telebot.SendOptions{
				DisableWebPagePreview: true,
				ParseMode:             telebot.ModeMarkdown,
			})
	} else {
		_, _ = B.Edit(msg, "订阅失败")
	}
}

//SendError send error user
func SendError(c *telebot.Chat) {
	B.Send(c, "请输入正确的指令！")
}

//BroadNews send new contents message to subscriber
func BroadNews(source *model.Source, subs []model.Subscribe, contents []model.Content) error {

	log.Printf("Source Title: <%s> Subscriber: %d New Contents: %d", source.Title, len(subs), len(contents))
	var u telebot.User
	for _, content := range contents {
		for _, sub := range subs {
			u.ID = int(sub.UserID)
			message := `
*%s*
%s
[原文](%s)
%s
`
			message = fmt.Sprintf(message, source.Title, content.Title, content.RawLink, content.TelegraphUrl)

			_, err := B.Send(&u, message, &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
