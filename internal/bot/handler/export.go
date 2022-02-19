package handler

import (
	tb "gopkg.in/telebot.v3"
)

type Export struct {
}

func (e *Export) Handle(ctx tb.Context) {
	//mention := message.MentionFromMessage(ctx.Message())
	//var sourceList []model.Source
	//var err error
	//if mention == "" {
	//	sourceList, err = model.GetSourcesByUserID(ctx.Chat().ID)
	//	if err != nil {
	//		zap.S().Warnf(err.Error())
	//		ctx.Send("导出失败")
	//		return
	//	}
	//} else {
	//	channelChat, err := ctx.Bot().ChatByID(mention)
	//	if err != nil {
	//		ctx.Send("导出失败")
	//		return
	//	}
	//
	//	adminList, err := ctx.Bot().AdminsOf(channelChat)
	//	if err != nil {
	//		ctx.Send("导出失败")
	//		return
	//	}
	//
	//	senderIsAdmin := false
	//	for _, admin := range adminList {
	//		if ctx.ChatMember().Sender.ID == admin.User.ID {
	//			senderIsAdmin = true
	//		}
	//	}
	//
	//	if !senderIsAdmin {
	//		ctx.Send(fmt.Sprintf("非频道管理员无法执行此操作"))
	//		return
	//	}
	//
	//	sourceList, err = model.GetSourcesByUserID(channelChat.ID)
	//	if err != nil {
	//		zap.S().Errorf(err.Error())
	//		ctx.Send("导出失败")
	//		return
	//	}
	//}
	//
	//if len(sourceList) == 0 {
	//	ctx.Send("订阅列表为空")
	//	return
	//}
	//
	//opmlStr, err := opml.ToOPML(sourceList)
	//if err != nil {
	//	ctx.Send("导出失败")
	//	return
	//}
	//opmlFile := &tb.Document{File: tb.FromReader(strings.NewReader(opmlStr))}
	//opmlFile.FileName = fmt.Sprintf("subscriptions_%d.opml", time.Now().Unix())
	//if err := ctx.Send(opmlFile); err != nil {
	//	ctx.Send("导出失败")
	//	zap.S().Errorf("send opml file failed, err:%+v", err)
	//}
}
