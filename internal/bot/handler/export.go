package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/opml"
)

type Export struct {
	core *core.Core
}

func NewExport(core *core.Core) *Export {
	return &Export{core: core}
}

func (e *Export) Description() string {
	return "导出OPML"
}

func (e *Export) Command() string {
	return "/export"
}

func (e *Export) getChannelSources(bot *tb.Bot, opUserID int64, channelName string) ([]*model.Source, error) {
	// 导出channel订阅
	channelChat, err := bot.ChatByUsername(channelName)
	if err != nil {
		return nil, errors.New("无法获取频道信息")
	}

	adminList, err := bot.AdminsOf(channelChat)
	if err != nil {
		return nil, errors.New("无法获取频道管理员信息")
	}

	senderIsAdmin := false
	for _, admin := range adminList {
		if opUserID == admin.User.ID {
			senderIsAdmin = true
			break
		}
	}

	if !senderIsAdmin {
		return nil, errors.New("非频道管理员无法执行此操作")
	}

	sources, err := e.core.GetUserSubscribedSources(context.Background(), channelChat.ID)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("获取订阅源信息失败")
	}
	return sources, nil
}

func (e *Export) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	var sources []*model.Source
	if mention == "" {
		var err error
		sources, err = e.core.GetUserSubscribedSources(context.Background(), ctx.Chat().ID)
		if err != nil {
			log.Error(err)
			return ctx.Send("导出失败")
		}
	} else {
		var err error
		sources, err = e.getChannelSources(ctx.Bot(), ctx.Chat().ID, mention)
		if err != nil {
			log.Error(err)
			return ctx.Send("导出失败")
		}
	}

	if len(sources) == 0 {
		return ctx.Send("订阅列表为空")
	}

	opmlStr, err := opml.ToOPML(sources)
	if err != nil {
		return ctx.Send("导出失败")
	}
	opmlFile := &tb.Document{File: tb.FromReader(strings.NewReader(opmlStr))}
	opmlFile.FileName = fmt.Sprintf("subscriptions_%d.opml", time.Now().Unix())
	if err := ctx.Send(opmlFile); err != nil {
		log.Errorf("send OPML file failed, err:%v", err)
		return ctx.Send("导出失败")
	}
	return nil
}

func (e *Export) Middlewares() []tb.MiddlewareFunc {
	return nil
}
