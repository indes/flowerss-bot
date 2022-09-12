package handler

import (
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/log"
)

type Start struct {
}

func NewStart() *Start {
	return &Start{}
}

func (s *Start) Command() string {
	return "/start"
}

func (s *Start) Description() string {
	return "开始使用"
}

func (s *Start) Handle(ctx tb.Context) error {
	log.Infof("/start id: %d", ctx.Chat().ID)
	return ctx.Send(fmt.Sprintf("你好，欢迎使用flowerss。"))
}

func (s *Start) Middlewares() []tb.MiddlewareFunc {
	return nil
}
