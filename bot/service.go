package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
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
		newText := fmt.Sprintf("频道 [%s](https://t.me/%s) 订阅 [%s](%s) 成功", channelChat.Title, channelChat.Username, source.Title, source.Link)
		_, err = B.Edit(msg, newText,
			&tb.SendOptions{
				DisableWebPagePreview: true,
				ParseMode:             tb.ModeMarkdown,
			})
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
	for _, content := range contents {

		previewText := trimDescription(content.Description, config.PreviewText)

		for _, sub := range subs {
			tpldata := &config.TplData{
				SourceTitle:     source.Title,
				ContentTitle:    content.Title,
				RawLink:         content.RawLink,
				PreviewText:     previewText,
				TelegraphURL:    content.TelegraphUrl,
				Tags:            sub.Tag,
				EnableTelegraph: sub.EnableTelegraph == 1 && content.TelegraphUrl != "",
			}

			u := &tb.User{
				ID: int(sub.UserID),
			}
			o := &tb.SendOptions{
				DisableWebPagePreview: config.DisableWebPagePreview,
				ParseMode:             config.MessageMode,
				DisableNotification:   sub.EnableNotification != 1,
			}
			msg, err := tpldata.Render(config.MessageMode)
			if err != nil {
				log.Println("BroadNews tpldata.Render err ", err)
				return
			}
			if _, err := B.Send(u, msg, o); err != nil {
				log.Println(err)
				if strings.Contains(err.Error(), "Forbidden") {
					log.Printf("Unsubscribe UserID:%d SourceID:%d", sub.UserID, sub.SourceID)
					sub.Unsub()
				}

				/*
					Telegram return error if markdown message has incomplete format.
					Print the msg to warn the user
					api error: Bad Request: can't parse entities: Can't find end of the entity starting at byte offset 894
				*/
				if strings.Contains(err.Error(), "parse entities") {
					log.Println("Markdown Err: ", msg)
				}
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

func CheckAdmin(upd *tb.Update) bool {

	if upd.Message != nil {
		if HasAdminType(upd.Message.Chat.Type) {
			adminList, _ := B.AdminsOf(upd.Message.Chat)
			for _, admin := range adminList {
				if admin.User.ID == upd.Message.Sender.ID {
					return true
				}
			}
			return false
		}

		return true
	} else if upd.Callback != nil {
		if HasAdminType(upd.Callback.Message.Chat.Type) {
			adminList, _ := B.AdminsOf(upd.Callback.Message.Chat)
			for _, admin := range adminList {
				if admin.User.ID == upd.Callback.Sender.ID {
					return true
				}
			}
			return false
		}

		return true
	}
	return false
}

func userIsAdminOfGroup(userID int, groupChat *tb.Chat) (isAdmin bool) {

	adminList, err := B.AdminsOf(groupChat)
	isAdmin = false

	if err != nil {
		return
	}

	for _, admin := range adminList {
		if userID == admin.User.ID {
			isAdmin = true
		}
	}
	return
}

func UserIsAdminChannel(userID int, channelChat *tb.Chat) (isAdmin bool) {
	adminList, err := B.AdminsOf(channelChat)
	isAdmin = false

	if err != nil {
		return
	}

	for _, admin := range adminList {
		if userID == admin.User.ID {
			isAdmin = true
		}
	}
	return
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

func GetMentionFromMessage(m *tb.Message) (mention string) {
	for _, entity := range m.Entities {
		if entity.Type == tb.EntityMention {
			if mention == "" {
				mention = m.Text[entity.Offset : entity.Offset+entity.Length]

			}
		}
	}
	return
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
