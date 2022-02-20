package bot

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

//func toggleCtrlButtons(c *tb.Callback, action string) {
//
//	if (c.Message.Chat.Type == tb.ChatGroup || c.Message.Chat.Type == tb.ChatSuperGroup) &&
//		!userIsAdminOfGroup(c.Sender.ID, c.Message.Chat) {
//		// check admin
//		return
//	}
//
//	data := strings.Split(c.Data, ":")
//	subscriberID, _ := strconv.Atoi(data[0])
//	// 如果订阅者与按钮点击者id不一致，需要验证管理员权限
//	if subscriberID != c.Sender.ID {
//		channelChat, err := B.ChatByID(fmt.Sprintf("%d", subscriberID))
//
//		if err != nil {
//			return
//		}
//
//		if !UserIsAdminChannel(c.Sender.ID, channelChat) {
//			return
//		}
//	}
//
//	msg := strings.Split(c.Message.Text, "\n")
//	subID, err := strconv.Atoi(strings.Split(msg[1], " ")[1])
//	if err != nil {
//		_ = B.Respond(c, &tb.CallbackResponse{
//			Text: "error",
//		})
//		return
//	}
//	sub, err := model.GetSubscribeByID(subID)
//	if sub == nil || err != nil {
//		_ = B.Respond(c, &tb.CallbackResponse{
//			Text: "error",
//		})
//		return
//	}
//
//	source, _ := model.GetSourceById(sub.SourceID)
//	t := template.New("setting template")
//	_, _ = t.Parse(feedSettingTmpl)
//
//	switch action {
//	case "toggleNotice":
//		err = sub.ToggleNotification()
//	case "toggleTelegraph":
//		err = sub.ToggleTelegraph()
//	case "toggleUpdate":
//		err = source.ToggleEnabled()
//	}
//
//	if err != nil {
//		_ = B.Respond(c, &tb.CallbackResponse{
//			Text: "error",
//		})
//		return
//	}
//
//	sub.Save()
//
//	text := new(bytes.Buffer)
//
//	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})
//	_ = B.Respond(c, &tb.CallbackResponse{
//		Text: "修改成功",
//	})
//	_, _ = B.Edit(c.Message, text.String(), &tb.SendOptions{
//		ParseMode: tb.ModeHTML,
//	}, &tb.ReplyMarkup{
//		InlineKeyboard: genFeedSetBtn(c, sub, source),
//	})
//}
//

//func setCmdCtr(m *tb.Message) {
//
//	mention := GetMentionFromMessage(m)
//	var sources []model.Source
//	var ownerID int64
//	// 获取订阅列表
//	if mention == "" {
//		sources, _ = model.GetSourcesByUserID(m.Chat.ID)
//		ownerID = int64(m.Chat.ID)
//		if len(sources) <= 0 {
//			_, _ = B.Send(m.Chat, "当前没有订阅源")
//			return
//		}
//
//	} else {
//
//		channelChat, err := B.ChatByID(mention)
//
//		if err != nil {
//			_, _ = B.Send(m.Chat, "获取Channel信息错误。")
//			return
//		}
//
//		if UserIsAdminChannel(m.Sender.ID, channelChat) {
//			sources, _ = model.GetSourcesByUserID(channelChat.ID)
//
//			if len(sources) <= 0 {
//				_, _ = B.Send(m.Chat, "Channel没有订阅源。")
//				return
//			}
//			ownerID = channelChat.ID
//
//		} else {
//			_, _ = B.Send(m.Chat, "非Channel管理员无法执行此操作。")
//			return
//		}
//
//	}
//
//	var replyButton []tb.ReplyButton
//	replyKeys := [][]tb.ReplyButton{}
//	setFeedItemBtns := [][]tb.InlineButton{}
//
//	// 配置按钮
//	for _, source := range sources {
//		// 添加按钮
//		text := fmt.Sprintf("%s %s", source.Title, source.Link)
//		replyButton = []tb.ReplyButton{
//			tb.ReplyButton{Text: text},
//		}
//		replyKeys = append(replyKeys, replyButton)
//
//		setFeedItemBtns = append(setFeedItemBtns, []tb.InlineButton{
//			tb.InlineButton{
//				Unique: "set_feed_item_btn",
//				Text:   fmt.Sprintf("[%d] %s", source.ID, source.Title),
//				Data:   fmt.Sprintf("%d:%d", ownerID, source.ID),
//			},
//		})
//	}
//
//	_, _ = B.Send(m.Chat, "请选择你要设置的源", &tb.ReplyMarkup{
//		InlineKeyboard: setFeedItemBtns,
//	})
//}
//
//func setFeedItemBtnCtr(c *tb.Callback) {
//
//	if (c.Message.Chat.Type == tb.ChatGroup || c.Message.Chat.Type == tb.ChatSuperGroup) &&
//		!userIsAdminOfGroup(c.Sender.ID, c.Message.Chat) {
//		return
//	}
//
//	data := strings.Split(c.Data, ":")
//	subscriberID, _ := strconv.Atoi(data[0])
//
//	// 如果订阅者与按钮点击者id不一致，需要验证管理员权限
//
//	if subscriberID != c.Sender.ID {
//		channelChat, err := B.ChatByID(fmt.Sprintf("%d", subscriberID))
//
//		if err != nil {
//			return
//		}
//
//		if !UserIsAdminChannel(c.Sender.ID, channelChat) {
//			return
//		}
//	}
//
//	sourceID, _ := strconv.Atoi(data[1])
//
//	source, err := model.GetSourceById(uint(sourceID))
//
//	if err != nil {
//		_, _ = B.Edit(c.Message, "找不到该订阅源，错误代码01。")
//		return
//	}
//
//	sub, err := model.GetSubscribeByUserIDAndSourceID(int64(subscriberID), source.ID)
//	if err != nil {
//		_, _ = B.Edit(c.Message, "用户未订阅该rss，错误代码02。")
//		return
//	}
//
//	t := template.New("setting template")
//	_, _ = t.Parse(feedSettingTmpl)
//	text := new(bytes.Buffer)
//	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})
//
//	_, _ = B.Edit(
//		c.Message,
//		text.String(),
//		&tb.SendOptions{
//			ParseMode: tb.ModeHTML,
//		}, &tb.ReplyMarkup{
//			InlineKeyboard: genFeedSetBtn(c, sub, source),
//		},
//	)
//}
//
//func setSubTagBtnCtr(c *tb.Callback) {
//
//	// 权限验证
//	if !feedSetAuth(c) {
//		return
//	}
//	data := strings.Split(c.Data, ":")
//	ownID, _ := strconv.Atoi(data[0])
//	sourceID, _ := strconv.Atoi(data[1])
//
//	sub, err := model.GetSubscribeByUserIDAndSourceID(int64(ownID), uint(sourceID))
//	if err != nil {
//		_, _ = B.Send(
//			c.Message.Chat,
//			"系统错误，代码04",
//		)
//		return
//	}
//	msg := fmt.Sprintf(
//		"请使用`/setfeedtag %d tags`命令为该订阅设置标签，tags为需要设置的标签，以空格分隔。（最多设置三个标签） \n"+
//			"例如：`/setfeedtag %d 科技 苹果`",
//		sub.ID, sub.ID)
//
//	_ = B.Delete(c.Message)
//
//	_, _ = B.Send(
//		c.Message.Chat,
//		msg,
//		&tb.SendOptions{ParseMode: tb.ModeMarkdown},
//	)
//}
//
//func genFeedSetBtn(c *tb.Callback, sub *model.Subscribe, source *model.Source) [][]tb.InlineButton {
//	setSubTagKey := tb.InlineButton{
//		Unique: "set_set_sub_tag_btn",
//		Text:   "标签设置",
//		Data:   c.Data,
//	}
//
//	toggleNoticeKey := tb.InlineButton{
//		Unique: "set_toggle_notice_btn",
//		Text:   "开启通知",
//		Data:   c.Data,
//	}
//	if sub.EnableNotification == 1 {
//		toggleNoticeKey.Text = "关闭通知"
//	}
//
//	toggleTelegraphKey := tb.InlineButton{
//		Unique: "set_toggle_telegraph_btn",
//		Text:   "开启 Telegraph 转码",
//		Data:   c.Data,
//	}
//	if sub.EnableTelegraph == 1 {
//		toggleTelegraphKey.Text = "关闭 Telegraph 转码"
//	}
//
//	toggleEnabledKey := tb.InlineButton{
//		Unique: "set_toggle_update_btn",
//		Text:   "暂停更新",
//		Data:   c.Data,
//	}
//
//	if source.ErrorCount >= config.ErrorThreshold {
//		toggleEnabledKey.Text = "重启更新"
//	}
//
//	feedSettingKeys := [][]tb.InlineButton{
//		[]tb.InlineButton{
//			toggleEnabledKey,
//			toggleNoticeKey,
//		},
//		[]tb.InlineButton{
//			toggleTelegraphKey,
//			setSubTagKey,
//		},
//	}
//	return feedSettingKeys
//}
//
//func setToggleNoticeBtnCtr(c *tb.Callback) {
//	toggleCtrlButtons(c, "toggleNotice")
//}
//
//func setToggleTelegraphBtnCtr(c *tb.Callback) {
//	toggleCtrlButtons(c, "toggleTelegraph")
//}
//
//func setToggleUpdateBtnCtr(c *tb.Callback) {
//	toggleCtrlButtons(c, "toggleUpdate")
//}

//func unsubFeedItemBtnCtr(c *tb.Callback) {
//	if (c.Message.Chat.Type == tb.ChatGroup || c.Message.Chat.Type == tb.ChatSuperGroup) &&
//		!userIsAdminOfGroup(c.Sender.ID, c.Message.Chat) {
//		// check admin
//		return
//	}
//
//	data := strings.Split(c.Data, ":")
//	if len(data) == 3 {
//		userID, _ := strconv.Atoi(data[0])
//		subID, _ := strconv.Atoi(data[1])
//		sourceID, _ := strconv.Atoi(data[2])
//		source, _ := model.GetSourceById(uint(sourceID))
//
//		rtnMsg := fmt.Sprintf("[%d] <a href=\"%s\">%s</a> 退订成功", sourceID, source.Link, source.Title)
//
//		err := model.UnsubByUserIDAndSubID(int64(userID), uint(subID))
//
//		if err == nil {
//			_, _ = B.Edit(
//				c.Message,
//				rtnMsg,
//				&tb.SendOptions{
//					ParseMode: tb.ModeHTML,
//				},
//			)
//			return
//		}
//	}
//	_, _ = B.Edit(c.Message, "退订错误！")
//}
//

//
//func setIntervalCmdCtr(m *tb.Message) {
//
//	args := strings.Split(m.Payload, " ")
//
//	if len(args) < 1 {
//		_, _ = B.Send(m.Chat, "/setinterval [interval] [sub id] 设置订阅刷新频率（可设置多个sub id，以空格分割）")
//		return
//	}
//
//	interval, err := strconv.Atoi(args[0])
//	if interval <= 0 || err != nil {
//		_, _ = B.Send(m.Chat, "请输入正确的抓取频率")
//		return
//	}
//
//	for _, id := range args[1:] {
//
//		subID, err := strconv.Atoi(id)
//		if err != nil {
//			_, _ = B.Send(m.Chat, "请输入正确的订阅id!")
//			return
//		}
//
//		sub, err := model.GetSubscribeByID(subID)
//
//		if err != nil || sub == nil {
//			_, _ = B.Send(m.Chat, "请输入正确的订阅id!")
//			return
//		}
//
//		if !checkPermit(int64(m.Sender.ID), sub.UserID) {
//			_, _ = B.Send(m.Chat, "没有权限!")
//			return
//		}
//
//		_ = sub.SetInterval(interval)
//
//	}
//	_, _ = B.Send(m.Chat, "抓取频率设置成功!")
//
//	return
//}

//func textCtr(m *tb.Message) {
//	switch UserState[m.Chat.ID] {
//	case fsm.UnSub:
//		{
//			str := strings.Split(m.Text, " ")
//
//			if len(str) < 2 && (strings.HasPrefix(str[0], "[") && strings.HasSuffix(str[0], "]")) {
//				_, _ = B.Send(m.Chat, "请选择正确的指令！")
//			} else {
//
//				var sourceID uint
//				if _, err := fmt.Sscanf(str[0], "[%d]", &sourceID); err != nil {
//					_, _ = B.Send(m.Chat, "请选择正确的指令！")
//					return
//				}
//
//				source, err := model.GetSourceById(sourceID)
//
//				if err != nil {
//					_, _ = B.Send(m.Chat, "请选择正确的指令！")
//					return
//				}
//
//				err = model.UnsubByUserIDAndSource(m.Chat.ID, source)
//
//				if err != nil {
//					_, _ = B.Send(m.Chat, "请选择正确的指令！")
//					return
//				}
//
//				_, _ = B.Send(
//					m.Chat,
//					fmt.Sprintf("[%s](%s) 退订成功", source.Title, source.Link),
//					&tb.SendOptions{
//						ParseMode: tb.ModeMarkdown,
//					}, &tb.ReplyMarkup{
//						ReplyKeyboardRemove: true,
//					},
//				)
//				UserState[m.Chat.ID] = fsm.None
//				return
//			}
//		}
//
//	case fsm.Sub:
//		{
//			url := strings.Split(m.Text, " ")
//			if !CheckURL(url[0]) {
//				_, _ = B.Send(m.Chat, "请回复正确的URL", &tb.ReplyMarkup{ForceReply: true})
//				return
//			}
//
//			registFeed(m.Chat, url[0])
//			UserState[m.Chat.ID] = fsm.None
//		}
//	case fsm.SetSubTag:
//		{
//			return
//		}
//	case fsm.Set:
//		{
//			str := strings.Split(m.Text, " ")
//			url := str[len(str)-1]
//			if len(str) != 2 && !CheckURL(url) {
//				_, _ = B.Send(m.Chat, "请选择正确的指令！")
//			} else {
//				source, err := model.GetSourceByUrl(url)
//
//				if err != nil {
//					_, _ = B.Send(m.Chat, "请选择正确的指令！")
//					return
//				}
//				sub, err := model.GetSubscribeByUserIDAndSourceID(m.Chat.ID, source.ID)
//				if err != nil {
//					_, _ = B.Send(m.Chat, "请选择正确的指令！")
//					return
//				}
//				t := template.New("setting template")
//				_, _ = t.Parse(feedSettingTmpl)
//
//				toggleNoticeKey := tb.InlineButton{
//					Unique: "set_toggle_notice_btn",
//					Text:   "开启通知",
//				}
//				if sub.EnableNotification == 1 {
//					toggleNoticeKey.Text = "关闭通知"
//				}
//
//				toggleTelegraphKey := tb.InlineButton{
//					Unique: "set_toggle_telegraph_btn",
//					Text:   "开启 Telegraph 转码",
//				}
//				if sub.EnableTelegraph == 1 {
//					toggleTelegraphKey.Text = "关闭 Telegraph 转码"
//				}
//
//				toggleEnabledKey := tb.InlineButton{
//					Unique: "set_toggle_update_btn",
//					Text:   "暂停更新",
//				}
//
//				if source.ErrorCount >= config.ErrorThreshold {
//					toggleEnabledKey.Text = "重启更新"
//				}
//
//				feedSettingKeys := [][]tb.InlineButton{
//					[]tb.InlineButton{
//						toggleEnabledKey,
//						toggleNoticeKey,
//						toggleTelegraphKey,
//					},
//				}
//
//				text := new(bytes.Buffer)
//
//				_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})
//
//				// send null message to remove old keyboard
//				delKeyMessage, err := B.Send(m.Chat, "processing", &tb.ReplyMarkup{ReplyKeyboardRemove: true})
//				err = B.Delete(delKeyMessage)
//
//				_, _ = B.Send(
//					m.Chat,
//					text.String(),
//					&tb.SendOptions{
//						ParseMode: tb.ModeHTML,
//					}, &tb.ReplyMarkup{
//						InlineKeyboard: feedSettingKeys,
//					},
//				)
//				UserState[m.Chat.ID] = fsm.None
//			}
//		}
//	}
//}
//
