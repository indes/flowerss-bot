package bot

import (
	"fmt"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

func FeedForChannelRegister(m *tb.Message, url string, channelMention string) {
	msg, err := B.Send(m.Chat, "处理中...")
	channelChat, err := B.ChatByID(channelMention)
	adminList, err := B.AdminsOf(channelChat)

	senderIsAdmin := false
	botIsAdmin := false
	for _, admin := range adminList {
		if m.Sender.ID == admin.User.ID {
			senderIsAdmin = true
		}
		if B.Me.ID == admin.User.ID {
			botIsAdmin = true
		}
	}

	if !botIsAdmin {
		msg, _ = B.Edit(msg, fmt.Sprintf("请先将bot添加为频道管理员"))
		return
	}

	if !senderIsAdmin {
		msg, _ = B.Edit(msg, fmt.Sprintf("非频道管理员无法执行此操作"))
		return
	}

	source, err := model.FindOrNewSourceByUrl(url)

	if err != nil {
		msg, _ = B.Edit(msg, fmt.Sprintf("%s，订阅失败", err))
		return
	}

	err = model.RegistFeed(channelChat.ID, source.ID)
	log.Printf("%d for %d subscribe [%d]%s %s", m.Chat.ID, channelChat.ID, source.ID, source.Title, source.Link)

	if err == nil {
		newText := fmt.Sprintf("成功为频道 [%s](https://t.me/%s) 订阅 [%s](%s) ", channelChat.Title, channelChat.Username, source.Title, source.Link)
		_, err = B.Edit(msg, newText,
			&tb.SendOptions{
				DisableWebPagePreview: true,
				ParseMode:             tb.ModeMarkdown,
			})
		log.Println(err)

	} else {
		_, _ = B.Edit(msg, "订阅失败")
	}
}

func registFeed(chat *tb.Chat, url string) {
	msg, err := B.Send(chat, "处理中...")

	source, err := model.FindOrNewSourceByUrl(url)

	if err != nil {
		msg, _ = B.Edit(msg, fmt.Sprintf("%s，订阅失败", err))
		return
	}

	err = model.RegistFeed(chat.ID, source.ID)
	log.Printf("%d subscribe [%d]%s %s", chat.ID, source.ID, source.Title, source.Link)

	if err == nil {
		_, _ = B.Edit(msg, fmt.Sprintf("[%s](%s) 订阅成功", source.Title, source.Link),
			&tb.SendOptions{
				DisableWebPagePreview: true,
				ParseMode:             tb.ModeMarkdown,
			})
	} else {
		_, _ = B.Edit(msg, "订阅失败")
	}
}

//SendError send error user
func SendError(c *tb.Chat) {
	_, _ = B.Send(c, "请输入正确的指令！")
}

//BroadNews send new contents message to subscriber
func BroadNews(source *model.Source, subs []model.Subscribe, contents []model.Content) {

	log.Printf("Source Title: <%s> Subscriber: %d New Contents: %d", source.Title, len(subs), len(contents))
	var u tb.User
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

			_, err := B.Send(&u, message, &tb.SendOptions{
				DisableWebPagePreview: false,
				ParseMode:             tb.ModeMarkdown,
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
	var u tb.User
	for _, sub := range subs {
		message := fmt.Sprintf("[%s](%s) 已经累计连续%d次更新失败，暂时停止更新", source.Title, source.Link, config.ErrorThreshold)
		u.ID = int(sub.UserID)
		_, _ = B.Send(&u, message, &tb.SendOptions{
			ParseMode: tb.ModeMarkdown,
		})
	}
}

func CheckAdmin(m *tb.Message) bool {
	if m == nil {
		return false
	}
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

func HasAdminType(t tb.ChatType) bool {
	hasAdmin := []tb.ChatType{tb.ChatGroup, tb.ChatSuperGroup, tb.ChatChannel, tb.ChatChannelPrivate}
	for _, n := range hasAdmin {
		if t == n {
			return true
		}
	}
	return false
}

func GetUrlAndMentionFromMessage(m *tb.Message) (url string, mention string) {
	for _, entity := range m.Entities {
		if entity.Type == tb.EntityMention {
			if mention == "" {
				mention = m.Text[entity.Offset : entity.Offset+entity.Length]

			}
		}

		if entity.Type == tb.EntityURL {
			if url == "" {
				url = m.Text[entity.Offset : entity.Offset+entity.Length]
			}
		}
	}

	return
}
