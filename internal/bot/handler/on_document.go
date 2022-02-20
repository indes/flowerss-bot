package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/opml"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
)

type OnDocument struct {
	bot *tb.Bot
}

func NewOnDocument(bot *tb.Bot) *OnDocument {
	return &OnDocument{bot: bot}
}

func (o *OnDocument) Command() string {
	return tb.OnDocument
}

func (o *OnDocument) Description() string {
	return ""
}

func (o *OnDocument) getOPML(ctx tb.Context) (*opml.OPML, error) {
	if !strings.HasSuffix(ctx.Message().Document.FileName, ".opml") {
		return nil, errors.New("请发送正确的 OPML 文件")
	}

	fileRead, err := o.bot.File(&ctx.Message().Document.File)
	if err != nil {
		return nil, errors.New("获取文件失败")
	}

	opmlFile, err := opml.ReadOPML(fileRead)
	if err != nil {
		zap.S().Errorf("parser opml failed, %v", err)
		return nil, errors.New("获取文件失败")
	}
	return opmlFile, nil
}

func (o *OnDocument) Handle(ctx tb.Context) error {
	opmlFile, err := o.getOPML(ctx)
	if err != nil {
		return ctx.Reply(err.Error())
	}
	userID := ctx.Chat().ID
	v := ctx.Get(session.StoreKeyMentionChat.String())
	if mentionChat, ok := v.(*tb.Chat); ok && mentionChat != nil {
		userID = mentionChat.ID
	}

	outlines, _ := opmlFile.GetFlattenOutlines()
	var failImportList []opml.Outline
	var successImportList []opml.Outline
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
		zap.S().Infof("%d subscribe [%d]%s %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
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

	return ctx.Reply(importReport, &tb.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tb.ModeHTML,
	})
}

func (o *OnDocument) Middlewares() []tb.MiddlewareFunc {
	return nil
}
