package bot

import (
	"bytes"
	"fmt"
	"github.com/indes/flowerss-bot/bot/fsm"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	feedSettingTmpl = `
订阅<b>设置</b>
[id] {{ .sub.ID }}
[标题] {{ .source.Title }}
[Link] {{.source.Link }}
[抓取更新] {{if ge .source.ErrorCount 100}}暂停{{else if lt .source.ErrorCount 100}}抓取中{{end}}
[通知] {{if eq .sub.EnableNotification 0}}关闭{{else if eq .sub.EnableNotification 1}}开启{{end}}
[Telegraph] {{if eq .sub.EnableTelegraph 0}}关闭{{else if eq .sub.EnableTelegraph 1}}开启{{end}}
[Tag] {{if .sub.Tag}}{{ .sub.Tag }}{{else}}无{{end}}
`
)

func toggleCtrlButtons(c *tb.Callback, action string) {

	if (c.Message.Chat.Type == tb.ChatGroup || c.Message.Chat.Type == tb.ChatSuperGroup) &&
		!userIsAdminOfGroup(c.Sender.ID, c.Message.Chat) {
		// check admin
		return
	}

	data := strings.Split(c.Data, ":")
	subscriberID, _ := strconv.Atoi(data[0])
	// 如果订阅者与按钮点击者id不一致，需要验证管理员权限
	if subscriberID != c.Sender.ID {
		channelChat, err := B.ChatByID(fmt.Sprintf("%d", subscriberID))

		if err != nil {
			return
		}

		if !UserIsAdminChannel(c.Sender.ID, channelChat) {
			return
		}
	}

	msg := strings.Split(c.Message.Text, "\n")
	subID, err := strconv.Atoi(strings.Split(msg[1], " ")[1])
	if err != nil {
		_ = B.Respond(c, &tb.CallbackResponse{
			Text: "error",
		})
		return
	}
	sub, err := model.GetSubscribeByID(subID)
	if sub == nil || err != nil {
		_ = B.Respond(c, &tb.CallbackResponse{
			Text: "error",
		})
		return
	}

	source, _ := model.GetSourceById(sub.SourceID)
	t := template.New("setting template")
	_, _ = t.Parse(feedSettingTmpl)

	switch action {
	case "toggleNotice":
		err = sub.ToggleNotification()
	case "toggleTelegraph":
		err = sub.ToggleTelegraph()
	case "toggleUpdate":
		err = source.ToggleEnabled()
	}

	if err != nil {
		_ = B.Respond(c, &tb.CallbackResponse{
			Text: "error",
		})
		return
	}

	sub.Save()

	text := new(bytes.Buffer)

	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub})
	_ = B.Respond(c, &tb.CallbackResponse{
		Text: "修改成功",
	})
	_, _ = B.Edit(c.Message, text.String(), &tb.SendOptions{
		ParseMode: tb.ModeHTML,
	}, &tb.ReplyMarkup{
		InlineKeyboard: genFeedSetBtn(c, sub, source),
	})
}

func startCmdCtr(m *tb.Message) {
	user := model.FindOrInitUser(m.Chat.ID)
	log.Printf("/start %d", user.ID)
	_, _ = B.Send(m.Chat, fmt.Sprintf("你好，欢迎使用flowerss。"))
}

func subCmdCtr(m *tb.Message) {
	url, mention := GetUrlAndMentionFromMessage(m)

	if mention == "" {
		if url != "" {
			registFeed(m.Chat, url)
		} else {
			_, err := B.Send(m.Chat, "请回复RSS URL", &tb.ReplyMarkup{ForceReply: true})
			if err == nil {
				UserState[m.Chat.ID] = fsm.Sub
			}
		}
	} else {
		if url != "" {
			FeedForChannelRegister(m, url, mention)
		} else {
			_, _ = B.Send(m.Chat, "频道订阅请使用' /sub @ChannelID URL ' 命令")
		}
	}

}

func exportCmdCtr(m *tb.Message) {

	mention := GetMentionFromMessage(m)
	var sourceList []model.Source
	var err error
	if mention == "" {

		sourceList, err = model.GetSourcesByUserID(m.Chat.ID)
		if err != nil {
			log.Println(err.Error())
			_, _ = B.Send(m.Chat, fmt.Sprintf("导出失败"))
			return
		}
	} else {
		channelChat, err := B.ChatByID(mention)

		if err != nil {
			_, _ = B.Send(m.Chat, "error")
			return
		}

		adminList, err := B.AdminsOf(channelChat)
		if err != nil {
			_, _ = B.Send(m.Chat, "error")
			return
		}

		senderIsAdmin := false
		for _, admin := range adminList {
			if m.Sender.ID == admin.User.ID {
				senderIsAdmin = true
			}
		}

		if !senderIsAdmin {
			_, _ = B.Send(m.Chat, fmt.Sprintf("非频道管理员无法执行此操作"))
			return
		}

		sourceList, err = model.GetSourcesByUserID(channelChat.ID)
		if err != nil {
			log.Println(err.Error())
			_, _ = B.Send(m.Chat, fmt.Sprintf("导出失败"))
			return
		}
	}

	if len(sourceList) == 0 {
		_, _ = B.Send(m.Chat, fmt.Sprintf("订阅列表为空"))
		return
	}

	opmlStr, err := ToOPML(sourceList)

	if err != nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("导出失败"))
		return
	}
	opmlFile := &tb.Document{File: tb.FromReader(strings.NewReader(opmlStr))}
	opmlFile.FileName = fmt.Sprintf("subscriptions_%d.opml", time.Now().Unix())
	_, err = B.Send(m.Chat, opmlFile)

	if err != nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("导出失败"))
		log.Println("[export]", err)
	}

}

func listCmdCtr(m *tb.Message) {
	mention := GetMentionFromMessage(m)
	if mention != "" {
		channelChat, err := B.ChatByID(mention)
		if err != nil {
			_, _ = B.Send(m.Chat, "error")
			return
		}
		adminList, err := B.AdminsOf(channelChat)
		if err != nil {
			_, _ = B.Send(m.Chat, "error")
			return
		}

		senderIsAdmin := false
		for _, admin := range adminList {
			if m.Sender.ID == admin.User.ID {
				senderIsAdmin = true
			}
		}

		if !senderIsAdmin {
			_, _ = B.Send(m.Chat, fmt.Sprintf("非频道管理员无法执行此操作"))
			return
		}

		sources, _ := model.GetSourcesByUserID(channelChat.ID)
		message := fmt.Sprintf("频道 [%s](https://t.me/%s) 订阅列表：\n", channelChat.Title, channelChat.Username)
		if len(sources) == 0 {
			message = fmt.Sprintf("频道 [%s](https://t.me/%s) 订阅列表为空", channelChat.Title, channelChat.Username)
		} else {
			for index, source := range sources {
				message = message + fmt.Sprintf("[[%d]] [%s](%s)\n", index+1, source.Title, source.Link)
			}
		}

		_, _ = B.Send(m.Chat, message, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})

	} else {
		sources, _ := model.GetSourcesByUserID(m.Chat.ID)
		message := "当前订阅列表：\n"
		if len(sources) == 0 {
			message = "订阅列表为空"
		} else {
			for index, source := range sources {
				message = message + fmt.Sprintf("[[%d]] [%s](%s)\n", index+1, source.Title, source.Link)
			}
		}
		_, _ = B.Send(m.Chat, message, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})
	}

}

func setCmdCtr(m *tb.Message) {

	mention := GetMentionFromMessage(m)
	var sources []model.Source
	var ownerID int64
	// 获取订阅列表
	if mention == "" {
		sources, _ = model.GetSourcesByUserID(m.Chat.ID)
		ownerID = int64(m.Sender.ID)
		if len(sources) <= 0 {
			_, _ = B.Send(m.Chat, "当前没有订阅源")
			return
		}

	} else {

		channelChat, err := B.ChatByID(mention)

		if err != nil {
			_, _ = B.Send(m.Chat, "获取Channel信息错误。")
			return
		}

		if UserIsAdminChannel(m.Sender.ID, channelChat) {
			sources, _ = model.GetSourcesByUserID(channelChat.ID)

			if len(sources) <= 0 {
				_, _ = B.Send(m.Chat, "Channel没有订阅源。")
				return
			}
			ownerID = channelChat.ID

		} else {
			_, _ = B.Send(m.Chat, "非Channel管理员无法执行此操作。")
			return
		}

	}

	var replyButton []tb.ReplyButton
	replyKeys := [][]tb.ReplyButton{}
	setFeedItemBtns := [][]tb.InlineButton{}

	// 配置按钮
	for _, source := range sources {
		// 添加按钮
		text := fmt.Sprintf("%s %s", source.Title, source.Link)
		replyButton = []tb.ReplyButton{
			tb.ReplyButton{Text: text},
		}
		replyKeys = append(replyKeys, replyButton)

		setFeedItemBtns = append(setFeedItemBtns, []tb.InlineButton{
			tb.InlineButton{
				Unique: "set_feed_item_btn",
				Text:   fmt.Sprintf("[%d] %s", source.ID, source.Title),
				Data:   fmt.Sprintf("%d:%d", ownerID, source.ID),
			},
		})
	}

	_, _ = B.Send(m.Chat, "请选择你要设置的源", &tb.ReplyMarkup{
		InlineKeyboard: setFeedItemBtns,
	})
}

func setFeedItemBtnCtr(c *tb.Callback) {

	if (c.Message.Chat.Type == tb.ChatGroup || c.Message.Chat.Type == tb.ChatSuperGroup) &&
		!userIsAdminOfGroup(c.Sender.ID, c.Message.Chat) {
		return
	}

	data := strings.Split(c.Data, ":")
	subscriberID, _ := strconv.Atoi(data[0])

	// 如果订阅者与按钮点击者id不一致，需要验证管理员权限

	if subscriberID != c.Sender.ID {
		channelChat, err := B.ChatByID(fmt.Sprintf("%d", subscriberID))

		if err != nil {
			return
		}

		if !UserIsAdminChannel(c.Sender.ID, channelChat) {
			return
		}
	}

	sourceID, _ := strconv.Atoi(data[1])

	source, err := model.GetSourceById(uint(sourceID))

	if err != nil {
		_, _ = B.Edit(c.Message, "找不到该订阅源，错误代码01。")
		return
	}

	sub, err := model.GetSubscribeByUserIDAndSourceID(int64(subscriberID), source.ID)
	if err != nil {
		_, _ = B.Edit(c.Message, "用户未订阅该rss，错误代码02。")
		return
	}

	t := template.New("setting template")
	_, _ = t.Parse(feedSettingTmpl)
	text := new(bytes.Buffer)
	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub})

	_, _ = B.Edit(
		c.Message,
		text.String(),
		&tb.SendOptions{
			ParseMode: tb.ModeHTML,
		}, &tb.ReplyMarkup{
			InlineKeyboard: genFeedSetBtn(c, sub, source),
		},
	)
}

func setSubTagBtnCtr(c *tb.Callback) {

	// 权限验证
	if !feedSetAuth(c) {
		return
	}
	data := strings.Split(c.Data, ":")
	subID, _ := strconv.Atoi(data[1])
	msg := fmt.Sprintf(
		"请使用`/setfeedtag 3 tags`命令为该订阅设置标签，tags为需要设置的标签，以空格分隔。（最多设置三个标签） \n例如：`/setfeedtag %d 科技 苹果`。",
		subID, subID)

	_ = B.Delete(c.Message)

	_, _ = B.Send(
		c.Message.Chat,
		msg,
		//&tb.ReplyMarkup{
		//	ForceReply: true,
		//	Selective:  true,
		//},
		&tb.SendOptions{ParseMode: tb.ModeMarkdown},
	)

	//if err == nil {
	//	UserState[c.Message.Chat.ID] = fsm.SetSubTag
	//}

}

func genFeedSetBtn(c *tb.Callback, sub *model.Subscribe, source *model.Source) [][]tb.InlineButton {
	setSubTagKey := tb.InlineButton{
		Unique: "set_set_sub_tag_btn",
		Text:   "标签设置",
		Data:   c.Data,
	}

	toggleNoticeKey := tb.InlineButton{
		Unique: "set_toggle_notice_btn",
		Text:   "开启通知",
		Data:   c.Data,
	}
	if sub.EnableNotification == 1 {
		toggleNoticeKey.Text = "关闭通知"
	}

	toggleTelegraphKey := tb.InlineButton{
		Unique: "set_toggle_telegraph_btn",
		Text:   "开启 Telegraph 转码",
		Data:   c.Data,
	}
	if sub.EnableTelegraph == 1 {
		toggleTelegraphKey.Text = "关闭 Telegraph 转码"
	}

	toggleEnabledKey := tb.InlineButton{
		Unique: "set_toggle_update_btn",
		Text:   "暂停更新",
		Data:   c.Data,
	}

	if source.ErrorCount >= config.ErrorThreshold {
		toggleEnabledKey.Text = "重启更新"
	}

	feedSettingKeys := [][]tb.InlineButton{
		[]tb.InlineButton{
			toggleEnabledKey,
			toggleNoticeKey,
		},
		[]tb.InlineButton{
			toggleTelegraphKey,
			setSubTagKey,
		},
	}
	return feedSettingKeys
}

func setToggleNoticeBtnCtr(c *tb.Callback) {
	toggleCtrlButtons(c, "toggleNotice")
}

func setToggleTelegraphBtnCtr(c *tb.Callback) {
	toggleCtrlButtons(c, "toggleTelegraph")
}

func setToggleUpdateBtnCtr(c *tb.Callback) {
	toggleCtrlButtons(c, "toggleUpdate")
}

func unsubCmdCtr(m *tb.Message) {

	url, mention := GetUrlAndMentionFromMessage(m)

	if mention == "" {
		if url != "" {
			//Unsub by url
			source, _ := model.GetSourceByUrl(url)
			if source == nil {
				_, _ = B.Send(m.Chat, "未订阅该RSS源")
			} else {
				err := model.UnsubByUserIDAndSource(m.Chat.ID, source)
				if err == nil {
					_, _ = B.Send(
						m.Chat,
						fmt.Sprintf("[%s](%s) 退订成功！", source.Title, source.Link),
						&tb.SendOptions{
							DisableWebPagePreview: true,
							ParseMode:             tb.ModeMarkdown,
						},
					)
					log.Printf("%d unsubscribe [%d]%s %s", m.Chat.ID, source.ID, source.Title, source.Link)
				} else {
					_, err = B.Send(m.Chat, err.Error())
				}
			}
		} else {
			//Unsub by button

			subs, err := model.GetSubsByUserID(m.Chat.ID)

			if err != nil {
				errorCtr(m, "Bot错误，请联系管理员！错误代码01")
				return
			}

			if len(subs) > 0 {
				unsubFeedItemBtns := [][]tb.InlineButton{}

				for _, sub := range subs {

					source, err := model.GetSourceById(sub.SourceID)
					if err != nil {
						errorCtr(m, "Bot错误，请联系管理员！错误代码02")
						return
					}

					unsubFeedItemBtns = append(unsubFeedItemBtns, []tb.InlineButton{
						tb.InlineButton{
							Unique: "unsub_feed_item_btn",
							Text:   fmt.Sprintf("[%d] %s", sub.SourceID, source.Title),
							Data:   fmt.Sprintf("%d:%d:%d", sub.UserID, sub.ID, source.ID),
							Action: nil,
						},
					})
				}

				_, _ = B.Send(m.Chat, "请选择你要退订的源", &tb.ReplyMarkup{
					InlineKeyboard: unsubFeedItemBtns,
				})
			} else {
				_, _ = B.Send(m.Chat, "当前没有订阅源")
			}
		}
	} else {
		if url != "" {
			channelChat, err := B.ChatByID(mention)
			if err != nil {
				_, _ = B.Send(m.Chat, "error")
				return
			}
			adminList, err := B.AdminsOf(channelChat)
			if err != nil {
				_, _ = B.Send(m.Chat, "error")
				return
			}

			senderIsAdmin := false
			for _, admin := range adminList {
				if m.Sender.ID == admin.User.ID {
					senderIsAdmin = true
				}
			}

			if !senderIsAdmin {
				_, _ = B.Send(m.Chat, fmt.Sprintf("非频道管理员无法执行此操作"))
				return
			}

			source, _ := model.GetSourceByUrl(url)
			sub, err := model.GetSubByUserIDAndURL(channelChat.ID, url)

			if err != nil {
				if err.Error() == "record not found" {
					_, _ = B.Send(
						m.Chat,
						fmt.Sprintf("频道 [%s](https://t.me/%s) 未订阅该RSS源", channelChat.Title, channelChat.Username),
						&tb.SendOptions{
							DisableWebPagePreview: true,
							ParseMode:             tb.ModeMarkdown,
						},
					)

				} else {
					_, _ = B.Send(m.Chat, "退订失败")
				}
				return

			} else {

				err := sub.Unsub()
				if err == nil {
					_, _ = B.Send(
						m.Chat,
						fmt.Sprintf("频道 [%s](https://t.me/%s) 退订 [%s](%s) 成功", channelChat.Title, channelChat.Username, source.Title, source.Link),
						&tb.SendOptions{
							DisableWebPagePreview: true,
							ParseMode:             tb.ModeMarkdown,
						},
					)
					log.Printf("%d for [%s]%s unsubscribe [%d]%s %s", m.Chat.ID, source.ID, source.Title, source.Link)
				} else {
					_, err = B.Send(m.Chat, err.Error())
				}
				return
			}

		} else {
			_, _ = B.Send(m.Chat, "频道退订请使用' /unsub @ChannelID URL ' 命令")
		}
	}

}

func unsubFeedItemBtnCtr(c *tb.Callback) {

	if (c.Message.Chat.Type == tb.ChatGroup || c.Message.Chat.Type == tb.ChatSuperGroup) &&
		!userIsAdminOfGroup(c.Sender.ID, c.Message.Chat) {
		// check admin
		return
	}

	data := strings.Split(c.Data, ":")
	if len(data) == 3 {
		userID, _ := strconv.Atoi(data[0])
		subID, _ := strconv.Atoi(data[1])
		sourceID, _ := strconv.Atoi(data[2])
		source, _ := model.GetSourceById(uint(sourceID))

		rtnMsg := fmt.Sprintf("[%d] <a href=\"%s\">%s</a> 退订成功", sourceID, source.Link, source.Title)

		err := model.UnsubByUserIDAndSubID(int64(userID), uint(subID))

		if err == nil {
			_, _ = B.Edit(
				c.Message,
				rtnMsg,
				&tb.SendOptions{
					ParseMode: tb.ModeHTML,
				},
			)
			return
		}
	}
	_, _ = B.Edit(c.Message, "退订错误！")
}

func unsubAllCmdCtr(m *tb.Message) {
	mention := GetMentionFromMessage(m)
	confirmKeys := [][]tb.InlineButton{}
	confirmKeys = append(confirmKeys, []tb.InlineButton{
		tb.InlineButton{
			Unique: "unsub_all_confirm_btn",
			Text:   "确认",
		},
		tb.InlineButton{
			Unique: "unsub_all_cancel_btn",
			Text:   "取消",
		},
	})

	var msg string

	if mention == "" {
		msg = "是否退订当前用户的所有订阅？"
	} else {
		msg = fmt.Sprintf("%s 是否退订该 Channel 所有订阅？", mention)
	}

	_, _ = B.Send(
		m.Chat,
		msg,
		&tb.SendOptions{
			ParseMode: tb.ModeHTML,
		}, &tb.ReplyMarkup{
			InlineKeyboard: confirmKeys,
		},
	)
}

func unsubAllCancelBtnCtr(c *tb.Callback) {
	_, _ = B.Edit(c.Message, "操作取消")
}

func unsubAllConfirmBtnCtr(c *tb.Callback) {
	mention := GetMentionFromMessage(c.Message)
	var msg string
	if mention == "" {
		success, fail, err := model.UnsubAllByUserID(int64(c.Sender.ID))
		if err != nil {
			msg = "退订失败"
		} else {
			msg = fmt.Sprintf("退订成功：%d\n退订失败：%d", success, fail)
		}

	} else {
		channelChat, err := B.ChatByID(mention)

		if err != nil {
			_, _ = B.Edit(c.Message, "error")
			return
		}

		if UserIsAdminChannel(c.Sender.ID, channelChat) {
			success, fail, err := model.UnsubAllByUserID(channelChat.ID)
			if err != nil {
				msg = "退订失败"

			} else {
				msg = fmt.Sprintf("退订成功：%d\n退订失败：%d", success, fail)
			}

		} else {
			msg = "非频道管理员无法执行此操作"
		}
	}

	_, _ = B.Edit(c.Message, msg)
}

func pingCmdCtr(m *tb.Message) {
	_, _ = B.Send(m.Chat, "pong")
}

func helpCmdCtr(m *tb.Message) {
	message := `
命令：
/sub 订阅源
/unsub  取消订阅
/list 查看当前订阅源
/set 设置订阅
/setfeedtag 设置订阅标签
/help 帮助
/import 导入 OPML 文件
/export 导出 OPML 文件
/unsuball 取消所有订阅
详细使用方法请看：https://github.com/indes/flowerss-bot
`
	_, _ = B.Send(m.Chat, message)
}

func importCmdCtr(m *tb.Message) {
	message := `请直接发送OPML文件。`
	_, _ = B.Send(m.Chat, message)
}

func setFeedTagCmdCtr(m *tb.Message) {
	args := strings.Split(m.Payload, " ")

	if len(args) < 1 {
		_, _ = B.Send(m.Chat, "命令错误")
		return
	}

	// 截短参数
	if len(args) > 4 {
		args = args[:4]
	}

	subID, err := strconv.Atoi(args[0])
	if err != nil {
		_, _ = B.Send(m.Chat, "请输入正确的订阅id!")
		return
	}

	sub, err := model.GetSubscribeByID(subID)

	if err != nil || sub == nil {
		return
	}

	if !checkPermit(int64(m.Sender.ID), sub.UserID) {
		_, _ = B.Send(m.Chat, "没有权限!")
		return
	}

	_ = sub.SetTag(args[1:])

	_, _ = B.Send(m.Chat, "订阅标签设置成功!")

	return
}

func textCtr(m *tb.Message) {
	switch UserState[m.Chat.ID] {
	case fsm.UnSub:
		{
			str := strings.Split(m.Text, " ")

			if len(str) < 2 && (strings.HasPrefix(str[0], "[") && strings.HasSuffix(str[0], "]")) {
				_, _ = B.Send(m.Chat, "请选择正确的指令！")
			} else {

				var sourceId uint
				if _, err := fmt.Sscanf(str[0], "[%d]", &sourceId); err != nil {
					_, _ = B.Send(m.Chat, "请选择正确的指令！")
					return
				}

				source, err := model.GetSourceById(sourceId)

				if err != nil {
					_, _ = B.Send(m.Chat, "请选择正确的指令！")
					return
				}

				err = model.UnsubByUserIDAndSource(m.Chat.ID, source)

				if err != nil {
					_, _ = B.Send(m.Chat, "请选择正确的指令！")
					return
				} else {
					_, _ = B.Send(
						m.Chat,
						fmt.Sprintf("[%s](%s) 退订成功", source.Title, source.Link),
						&tb.SendOptions{
							ParseMode: tb.ModeMarkdown,
						}, &tb.ReplyMarkup{
							ReplyKeyboardRemove: true,
						},
					)
					UserState[m.Chat.ID] = fsm.None
					return
				}
			}
		}

	case fsm.Sub:
		{
			url := strings.Split(m.Text, " ")
			if !CheckUrl(url[0]) {
				_, _ = B.Send(m.Chat, "请回复正确的URL", &tb.ReplyMarkup{ForceReply: true})
				return
			}

			registFeed(m.Chat, url[0])
			UserState[m.Chat.ID] = fsm.None
		}
	case fsm.SetSubTag:
		{
			return
		}
	case fsm.Set:
		{

			str := strings.Split(m.Text, " ")
			url := str[len(str)-1]
			if len(str) != 2 && !CheckUrl(url) {
				_, _ = B.Send(m.Chat, "请选择正确的指令！")
			} else {
				source, err := model.GetSourceByUrl(url)

				if err != nil {
					_, _ = B.Send(m.Chat, "请选择正确的指令！")
					return
				}
				sub, err := model.GetSubscribeByUserIDAndSourceID(m.Chat.ID, source.ID)
				if err != nil {
					_, _ = B.Send(m.Chat, "请选择正确的指令！")
					return
				}
				t := template.New("setting template")
				_, _ = t.Parse(feedSettingTmpl)

				toggleNoticeKey := tb.InlineButton{
					Unique: "set_toggle_notice_btn",
					Text:   "开启通知",
				}
				if sub.EnableNotification == 1 {
					toggleNoticeKey.Text = "关闭通知"
				}

				toggleTelegraphKey := tb.InlineButton{
					Unique: "set_toggle_telegraph_btn",
					Text:   "开启 Telegraph 转码",
				}
				if sub.EnableTelegraph == 1 {
					toggleTelegraphKey.Text = "关闭 Telegraph 转码"
				}

				toggleEnabledKey := tb.InlineButton{
					Unique: "set_toggle_update_btn",
					Text:   "暂停更新",
				}

				if source.ErrorCount >= config.ErrorThreshold {
					toggleEnabledKey.Text = "重启更新"
				}

				feedSettingKeys := [][]tb.InlineButton{
					[]tb.InlineButton{
						toggleEnabledKey,
						toggleNoticeKey,
						toggleTelegraphKey,
					},
				}

				text := new(bytes.Buffer)

				_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub})

				// send null message to remove old keyboard
				delKeyMessage, err := B.Send(m.Chat, "processing", &tb.ReplyMarkup{ReplyKeyboardRemove: true})
				err = B.Delete(delKeyMessage)

				_, _ = B.Send(
					m.Chat,
					text.String(),
					&tb.SendOptions{
						ParseMode: tb.ModeHTML,
					}, &tb.ReplyMarkup{
						InlineKeyboard: feedSettingKeys,
					},
				)
				UserState[m.Chat.ID] = fsm.None
			}
		}
	}
}

func docCtr(m *tb.Message) {
	if m.FromGroup() {
		if !userIsAdminOfGroup(m.ID, m.Chat) {
			return
		}
	}

	if m.FromChannel() {
		if !UserIsAdminChannel(m.ID, m.Chat) {
			return
		}
	}

	url, _ := B.FileURLByID(m.Document.FileID)
	if !strings.HasSuffix(url, ".opml") {
		return
	}

	opml, err := GetOPMLByURL(url)
	if err != nil {
		if err.Error() == "fetch opml file error" {
			_, _ = B.Send(m.Chat,
				"下载 OPML 文件失败，请检查 bot 服务器能否正常连接至 Telegram 服务器或稍后尝试导入。错误代码 02")

		} else {
			_, _ = B.Send(
				m.Chat,
				fmt.Sprintf(
					"如果需要导入订阅，请发送正确的 OPML 文件。错误代码 01，doc mimetype: %s",
					m.Document.MIME),
			)
		}
		return
	}

	message, _ := B.Send(m.Chat, "处理中，请稍后...")
	outlines, _ := opml.GetFlattenOutlines()
	var failImportList []Outline
	var successImportList []Outline

	for _, outline := range outlines {
		source, err := model.FindOrNewSourceByUrl(outline.XMLURL)
		if err != nil {
			failImportList = append(failImportList, outline)
			continue
		}
		err = model.RegistFeed(m.Chat.ID, source.ID)
		if err != nil {
			failImportList = append(failImportList, outline)
			continue
		}
		log.Printf("%d subscribe [%d]%s %s", m.Chat.ID, source.ID, source.Title, source.Link)
		successImportList = append(successImportList, outline)
	}

	importReport := fmt.Sprintf("<b>导入成功：%d，导入失败：%d</b>", len(successImportList), len(failImportList))
	if len(successImportList) != 0 {
		successReport := "\n\n<b>以下订阅源导入成功:</b>"
		for i, line := range successImportList {
			if line.Text != "" {
				successReport += fmt.Sprintf("\n[%d] <a href=\"%s\">%s</a>", i+1, line.XMLURL, line.Text)
			} else {
				successReport += fmt.Sprintf("\n[%d] %s", i+1, line.XMLURL)
			}
		}
		importReport += successReport
	}

	if len(failImportList) != 0 {
		failReport := "\n\n<b>以下订阅源导入失败:</b>"
		for i, line := range failImportList {
			if line.Text != "" {
				failReport += fmt.Sprintf("\n[%d] <a href=\"%s\">%s</a>", i+1, line.XMLURL, line.Text)
			} else {
				failReport += fmt.Sprintf("\n[%d] %s", i+1, line.XMLURL)
			}
		}
		importReport += failReport
	}
	_, err = B.Edit(message, importReport, &tb.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tb.ModeHTML,
	})

	//if err != nil {
	//	log.Println(err.Error())
	//}
	//if m.Document.MIME == "text/x-opml+xml" || m.Document.MIME == "application/xml" {
	//
	//} else {
	//	_, _ = B.Send(m.Chat, "如果需要导入订阅，请发送正确的OPML文件。错误代码 01")
	//}

}

func errorCtr(m *tb.Message, errMsg string) {
	_, _ = B.Send(m.Chat, errMsg)
}
