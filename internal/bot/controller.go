package bot

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
