package handler

import (
	tb "gopkg.in/telebot.v3"
)

type Ping struct {
	bot *tb.Bot
}

// NewPing new ping cmd handler
func NewPing(bot *tb.Bot) *Ping {
	return &Ping{bot: bot}
}

func (p *Ping) Command() string {
	return "/ping"
}

func (p *Ping) Description() string {
	return ""
}

func (p *Ping) Handle(ctx tb.Context) error {
	return ctx.Send("pong")
}

func (p *Ping) Middlewares() []tb.MiddlewareFunc {
	return nil
}
