package bot

import (
	"fmt"
	"github.com/indes/rssflow/config"
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
	_, _ = B.Send(c, "请输入正确的指令！")
}

//BroadNews send new contents message to subscriber
func BroadNews(source *model.Source, subs []model.Subscribe, contents []model.Content) {

	log.Printf("Source Title: <%s> Subscriber: %d New Contents: %d", source.Title, len(subs), len(contents))
	var u telebot.User
	var message string
	for _, content := range contents {
		for _, sub := range subs {
			var disableNotification bool
			if sub.EnableNotification == 1 {
				disableNotification = false
			} else {
				disableNotification = true
			}

			u.ID = int(sub.UserID)

			if sub.EnableTelegraph == 1 && content.TelegraphUrl != "" {
				message = `
*%s*
%s | [Telegraph](%s) | [原文](%s)
`
				message = fmt.Sprintf(message, source.Title, content.Title, content.TelegraphUrl, content.RawLink)
			} else {
				message = `
*%s*
%s | [原文](%s)
`
				message = fmt.Sprintf(message, source.Title, content.Title, content.RawLink)
			}

			_, err := B.Send(&u, message, &telebot.SendOptions{
				DisableWebPagePreview: false,
				ParseMode:             telebot.ModeMarkdown,
				DisableNotification:   disableNotification,
			})
			if err != nil {
				log.Println(err)
			}

		}
	}
}

func BroadSourceError(source *model.Source) {
	subs := model.GetSubscriberBySource(source)
	var u telebot.User
	for _, sub := range subs {
		message := fmt.Sprintf("[%s](%s) 已经累计连续%d次更新失败，暂时停止更新", source.Title, source.Link, config.ErrorThreshold)
		u.ID = int(sub.UserID)
		_, _ = B.Send(&u, message, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	}
}

func CheckAdmin(m *telebot.Message) bool {
	if HasAdminType(m.Chat.Type) {
		adminList, _ := B.AdminsOf(m.Chat)
		for _, admin := range adminList {
			if admin.User.ID == m.Sender.ID {
				return true
			}
		}
		return false
	}
	return true
}

func HasAdminType(t telebot.ChatType) bool {
	hasAdmin := []telebot.ChatType{telebot.ChatGroup, telebot.ChatSuperGroup, telebot.ChatChannel, telebot.ChatChannelPrivate}
	for _, n := range hasAdmin {
		if t == n {
			return true
		}
	}
	return false
}
