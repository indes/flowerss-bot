package handler

import tb "gopkg.in/telebot.v3"

type Ping struct {
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
