package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/opml"

	tb "gopkg.in/telebot.v3"
)

type OnDocument struct {
	bot  *tb.Bot
	core *core.Core
}

func NewOnDocument(bot *tb.Bot, core *core.Core) *OnDocument {
	return &OnDocument{
		bot:  bot,
		core: core,
	}
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
		log.Errorf("parser opml failed, %v", err)
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
	var failImportList = make([]opml.Outline, len(outlines))
	failIndex := 0
	var successImportList = make([]opml.Outline, len(outlines))
	successIndex := 0
	wg := &sync.WaitGroup{}
	for _, outline := range outlines {
		outline := outline
		wg.Add(1)
		go func() {
			defer wg.Done()
			source, err := o.core.CreateSource(context.Background(), outline.XMLURL)
			if err != nil {
				failImportList[failIndex] = outline
				failIndex++
				return
			}

			err = o.core.AddSubscription(context.Background(), userID, source.ID)
			if err != nil {
				if err == core.ErrSubscriptionExist {
					successImportList[successIndex] = outline
					successIndex++
				} else {
					failImportList[failIndex] = outline
					failIndex++
				}
				return
			}

			log.Infof("%d subscribe [%d]%s %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
			successImportList[successIndex] = outline
			successIndex++
			return
		}()
	}
	wg.Wait()

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("<b>导入成功：%d，导入失败：%d</b>\n", successIndex, failIndex))
	if successIndex != 0 {
		msg.WriteString("<b>以下订阅源导入成功:</b>\n")
		for i := 0; i < successIndex; i++ {
			line := successImportList[i]
			if line.Text != "" {
				msg.WriteString(
					fmt.Sprintf("[%d] <a href=\"%s\">%s</a>\n", i+1, line.XMLURL, line.Text),
				)
			} else {
				msg.WriteString(fmt.Sprintf("[%d] %s\n", i+1, line.XMLURL))
			}
		}

		msg.WriteString("\n")
	}

	if failIndex != 0 {
		msg.WriteString("<b>以下订阅源导入失败:</b>\n")
		for i := 0; i < failIndex; i++ {
			line := failImportList[i]
			if line.Text != "" {
				msg.WriteString(fmt.Sprintf("[%d] <a href=\"%s\">%s</a>\n", i+1, line.XMLURL, line.Text))
			} else {
				msg.WriteString(fmt.Sprintf("[%d] %s\n", i+1, line.XMLURL))
			}
		}

	}

	return ctx.Reply(
		msg.String(), &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeHTML,
		},
	)
}

func (o *OnDocument) Middlewares() []tb.MiddlewareFunc {
	return nil
}
