package bot

import (
	"fmt"
	"github.com/indes/rssflow/model"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

func registFeed(chat *telebot.Chat, url string) {
	msg, _ := B.Send(chat, "处理中...")

	source, _ := model.FindOrNewSourceByUrl(url)
	err := model.RegistFeed(chat.ID, source.ID)
	log.Printf("%d subscribe %s %s", chat.ID, source.Title, source.Link)

	if err == nil {
		msg, _ = B.Edit(msg, fmt.Sprintf("<%s> 订阅成功", source.Title))
	} else {
		msg, _ = B.Edit(msg, "订阅失败")
	}
}

func SendError(c *telebot.Chat) {
	B.Send(c, "请输入正确的指令！")
}

func BroadNews(source *model.Source, subs []model.Subscribe, contents []model.Content) error {
	//log.Println("Subs Len: ", len(subs), " Contents Len: ", len(contents))
	log.Printf("Source Title: <%s> Subscriber: %d New Contents: %d", source.Title, len(subs), len(contents))
	var u telebot.User
	for _, content := range contents {
		for _, sub := range subs {
			u.ID = int(sub.UserID)
			message := `
*%s*
%s
[原文](%s)
[Telegraph](%s)
`
			message = fmt.Sprintf(message, source.Title, content.Title, content.RawLink, content.TelegraphUrl)

			_, err := B.Send(&u, message, &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
			log.Println(err)
		}
	}
	return nil
}
