package handler

import (
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/config"
)

type Version struct {
}

func NewVersion() *Version {
	return &Version{}
}

func (c *Version) Command() string {
	return "/version"
}

func (c *Version) Description() string {
	return "Bot 版本信息"
}

func (c *Version) Handle(ctx tb.Context) error {
	return ctx.Send(config.AppVersionInfo())
}

func (c *Version) Middlewares() []tb.MiddlewareFunc {
	return nil
}
