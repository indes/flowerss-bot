package bot

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/indes/flowerss-bot/bot/fsm"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/log"
	"github.com/indes/flowerss-bot/model"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	feedSettingTmpl = `
订阅<b>设置</b>
[id] {{ .sub.ID }}
[标题] {{ .source.Title }}
[Link] {{.source.Link }}
[抓取更新] {{if ge .source.ErrorCount .Count }}暂停{{else if lt .source.ErrorCount .Count }}抓取中{{end}}
[抓取频率] {{ .sub.Interval }}分钟
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

	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})
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
	user, _ := model.FindOrCreateUserByTelegramID(m.Chat.ID)
	log.Printf("/start user_id: %d telegram_id: %d", user.ID, user.TelegramID)
	_, _ = B.Send(m.Chat, fmt.Sprintf("你好，欢迎使用flowerss。"))
}

func subCmdCtr(m *tb.Message) {

	url, mention := GetURLAndMentionFromMessage(m)

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

	var rspMessage string
	if mention != "" {
		// channel feed list
		channelChat, err := B.ChatByID(mention)
		if err != nil {
			_, _ = B.Send(m.Chat, "error")
			return
		}

		if !checkPermitOfChat(int64(m.Sender.ID), channelChat) {
			B.Send(m.Chat, fmt.Sprintf("非频道管理员无法执行此操作"))
			return
		}

		user, err := model.FindOrCreateUserByTelegramID(channelChat.ID)
		if err != nil {
			B.Send(m.Chat, fmt.Sprintf("内部错误 list@1"))
			return
		}

		subSourceMap, err := user.GetSubSourceMap()
		if err != nil {
			B.Send(m.Chat, fmt.Sprintf("内部错误 list@2"))
			return
		}

		sources, _ := model.GetSourcesByUserID(channelChat.ID)
		rspMessage = fmt.Sprintf("频道 [%s](https://t.me/%s) 订阅列表：\n", channelChat.Title, channelChat.Username)
		if len(sources) == 0 {
			rspMessage = fmt.Sprintf("频道 [%s](https://t.me/%s) 订阅列表为空", channelChat.Title, channelChat.Username)
		} else {
			for sub, source := range subSourceMap {
				rspMessage = rspMessage + fmt.Sprintf("[[%d]] [%s](%s)\n", sub.ID, source.Title, source.Link)
			}
		}
	} else {
		// private chat or group
		if m.Chat.Type != tb.ChatPrivate && !checkPermitOfChat(int64(m.Sender.ID), m.Chat) {
			// 无权限
			return
		}

		user, err := model.FindOrCreateUserByTelegramID(m.Chat.ID)
		if err != nil {
			B.Send(m.Chat, fmt.Sprintf("内部错误 list@1"))
			return
		}

		subSourceMap, err := user.GetSubSourceMap()
		if err != nil {
			B.Send(m.Chat, fmt.Sprintf("内部错误 list@2"))
			return
		}

		rspMessage = "当前订阅列表：\n"
		if len(subSourceMap) == 0 {
			rspMessage = "订阅列表为空"
		} else {
			for sub, source := range subSourceMap {
				rspMessage = rspMessage + fmt.Sprintf("[[%d]] [%s](%s)\n", sub.ID, source.Title, source.Link)
			}
		}
	}
	_, _ = B.Send(m.Chat, rspMessage, &tb.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tb.ModeMarkdown,
	})
}

func checkCmdCtr(m *tb.Message) {
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

		sources, _ := model.GetErrorSourcesByUserID(channelChat.ID)
		message := fmt.Sprintf("频道 [%s](https://t.me/%s) 失效订阅的列表：\n", channelChat.Title, channelChat.Username)
		if len(sources) == 0 {
			message = fmt.Sprintf("频道 [%s](https://t.me/%s) 所有订阅正常", channelChat.Title, channelChat.Username)
		} else {
			for _, source := range sources {
				message = message + fmt.Sprintf("[[%d]] [%s](%s)\n", source.ID, source.Title, source.Link)
			}
		}

		_, _ = B.Send(m.Chat, message, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})

	} else {
		sources, _ := model.GetErrorSourcesByUserID(m.Chat.ID)
		message := "失效订阅的列表：\n"
		if len(sources) == 0 {
			message = "所有订阅正常"
		} else {
			for _, source := range sources {
				message = message + fmt.Sprintf("[[%d]] [%s](%s)\n", source.ID, source.Title, source.Link)
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
		ownerID = int64(m.Chat.ID)
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
	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})

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
	ownID, _ := strconv.Atoi(data[0])
	sourceID, _ := strconv.Atoi(data[1])

	sub, err := model.GetSubscribeByUserIDAndSourceID(int64(ownID), uint(sourceID))
	if err != nil {
		_, _ = B.Send(
			c.Message.Chat,
			"系统错误，代码04",
		)
		return
	}
	msg := fmt.Sprintf(
		"请使用`/setfeedtag %d tags`命令为该订阅设置标签，tags为需要设置的标签，以空格分隔。（最多设置三个标签） \n"+
			"例如：`/setfeedtag %d 科技 苹果`",
		sub.ID, sub.ID)

	_ = B.Delete(c.Message)

	_, _ = B.Send(
		c.Message.Chat,
		msg,
		&tb.SendOptions{ParseMode: tb.ModeMarkdown},
	)
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

	url, mention := GetURLAndMentionFromMessage(m)

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

			}

			err = sub.Unsub()
			if err == nil {
				_, _ = B.Send(
					m.Chat,
					fmt.Sprintf("频道 [%s](https://t.me/%s) 退订 [%s](%s) 成功", channelChat.Title, channelChat.Username, source.Title, source.Link),
					&tb.SendOptions{
						DisableWebPagePreview: true,
						ParseMode:             tb.ModeMarkdown,
					},
				)
				log.Printf("%d for [%d]%s unsubscribe %s", m.Chat.ID, source.ID, source.Title, source.Link)
			} else {
				_, err = B.Send(m.Chat, err.Error())
			}
			return

		}
		_, _ = B.Send(m.Chat, "频道退订请使用' /unsub @ChannelID URL ' 命令")
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
	log.DebugWithMessage(m, "ping")
}

func helpCmdCtr(m *tb.Message) {
	message := `
命令：
/sub 订阅源
/unsub  取消订阅
/list 查看当前订阅源
/set 设置订阅
/check 检查当前订阅
/setfeedtag 设置订阅标签
/setinterval 设置订阅刷新频率
/activeall 开启所有订阅
/pauseall 暂停所有订阅
/help 帮助
/import 导入 OPML 文件
/export 导出 OPML 文件
/unsuball 取消所有订阅
详细使用方法请看：https://github.com/indes/flowerss-bot
`

	_, _ = B.Send(m.Chat, message)
}

func versionCmdCtr(m *tb.Message) {
	_, _ = B.Send(m.Chat, config.AppVersionInfo())
}

func importCmdCtr(m *tb.Message) {
	message := `请直接发送OPML文件，
如果需要为channel导入OPML，请在发送文件的时候附上channel id，例如@telegram
`
	_, _ = B.Send(m.Chat, message)
}

func setFeedTagCmdCtr(m *tb.Message) {
	mention := GetMentionFromMessage(m)
	args := strings.Split(m.Payload, " ")

	if len(args) < 1 {
		B.Send(m.Chat, "/setfeedtag [sub id] [tag1] [tag2] 设置订阅标签（最多设置三个Tag，以空格分割）")
		return
	}

	var subID int
	var err error
	if mention == "" {
		// 截短参数
		if len(args) > 4 {
			args = args[:4]
		}
		subID, err = strconv.Atoi(args[0])
		if err != nil {
			B.Send(m.Chat, "请输入正确的订阅id!")
			return
		}
	} else {
		if len(args) > 5 {
			args = args[:5]
		}
		subID, err = strconv.Atoi(args[1])
		if err != nil {
			B.Send(m.Chat, "请输入正确的订阅id!")
			return
		}
	}

	sub, err := model.GetSubscribeByID(subID)
	if err != nil || sub == nil {
		B.Send(m.Chat, "请输入正确的订阅id!")
		return
	}

	if !checkPermit(int64(m.Sender.ID), sub.UserID) {
		B.Send(m.Chat, "没有权限!")
		return
	}

	if mention == "" {
		err = sub.SetTag(args[1:])
	} else {
		err = sub.SetTag(args[2:])
	}

	if err != nil {
		B.Send(m.Chat, "订阅标签设置失败!")
		return
	}
	B.Send(m.Chat, "订阅标签设置成功!")
}

func setIntervalCmdCtr(m *tb.Message) {

	args := strings.Split(m.Payload, " ")

	if len(args) < 1 {
		_, _ = B.Send(m.Chat, "/setinterval [interval] [sub id] 设置订阅刷新频率（可设置多个sub id，以空格分割）")
		return
	}

	interval, err := strconv.Atoi(args[0])
	if interval <= 0 || err != nil {
		_, _ = B.Send(m.Chat, "请输入正确的抓取频率")
		return
	}

	for _, id := range args[1:] {

		subID, err := strconv.Atoi(id)
		if err != nil {
			_, _ = B.Send(m.Chat, "请输入正确的订阅id!")
			return
		}

		sub, err := model.GetSubscribeByID(subID)

		if err != nil || sub == nil {
			_, _ = B.Send(m.Chat, "请输入正确的订阅id!")
			return
		}

		if !checkPermit(int64(m.Sender.ID), sub.UserID) {
			_, _ = B.Send(m.Chat, "没有权限!")
			return
		}

		_ = sub.SetInterval(interval)

	}
	_, _ = B.Send(m.Chat, "抓取频率设置成功!")

	return
}

func activeAllCmdCtr(m *tb.Message) {
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

		_ = model.ActiveSourcesByUserID(channelChat.ID)
		message := fmt.Sprintf("频道 [%s](https://t.me/%s) 订阅已全部开启", channelChat.Title, channelChat.Username)

		_, _ = B.Send(m.Chat, message, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})

	} else {
		_ = model.ActiveSourcesByUserID(m.Chat.ID)
		message := "订阅已全部开启"

		_, _ = B.Send(m.Chat, message, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})
	}

}

func pauseAllCmdCtr(m *tb.Message) {
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

		_ = model.PauseSourcesByUserID(channelChat.ID)
		message := fmt.Sprintf("频道 [%s](https://t.me/%s) 订阅已全部暂停", channelChat.Title, channelChat.Username)

		_, _ = B.Send(m.Chat, message, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})

	} else {
		_ = model.PauseSourcesByUserID(m.Chat.ID)
		message := "订阅已全部暂停"

		_, _ = B.Send(m.Chat, message, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})
	}

}

func textCtr(m *tb.Message) {
	switch UserState[m.Chat.ID] {
	case fsm.UnSub:
		{
			str := strings.Split(m.Text, " ")

			if len(str) < 2 && (strings.HasPrefix(str[0], "[") && strings.HasSuffix(str[0], "]")) {
				_, _ = B.Send(m.Chat, "请选择正确的指令！")
			} else {

				var sourceID uint
				if _, err := fmt.Sscanf(str[0], "[%d]", &sourceID); err != nil {
					_, _ = B.Send(m.Chat, "请选择正确的指令！")
					return
				}

				source, err := model.GetSourceById(sourceID)

				if err != nil {
					_, _ = B.Send(m.Chat, "请选择正确的指令！")
					return
				}

				err = model.UnsubByUserIDAndSource(m.Chat.ID, source)

				if err != nil {
					_, _ = B.Send(m.Chat, "请选择正确的指令！")
					return
				}

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

	case fsm.Sub:
		{
			url := strings.Split(m.Text, " ")
			if !CheckURL(url[0]) {
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
			if len(str) != 2 && !CheckURL(url) {
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

				_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})

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

// docCtr Document handler
func docCtr(m *tb.Message) {
	if m.FromGroup() {
		if !userIsAdminOfGroup(m.Sender.ID, m.Chat) {
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
		B.Send(m.Chat, "如果需要导入订阅，请发送正确的 OPML 文件。")
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

	userID := m.Chat.ID
	mention := GetMentionFromMessage(m)
	if mention != "" {
		// import for channel
		channelChat, err := B.ChatByID(mention)
		if err != nil {
			_, _ = B.Send(m.Chat, "获取channel信息错误，请检查channel id是否正确")
			return
		}

		if !checkPermitOfChat(int64(m.Sender.ID), channelChat) {
			_, _ = B.Send(m.Chat, fmt.Sprintf("非频道管理员无法执行此操作"))
			return
		}

		userID = channelChat.ID
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
		err = model.RegistFeed(userID, source.ID)
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

	_, _ = B.Edit(message, importReport, &tb.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tb.ModeHTML,
	})
}

func errorCtr(m *tb.Message, errMsg string) {
	_, _ = B.Send(m.Chat, errMsg)
}
