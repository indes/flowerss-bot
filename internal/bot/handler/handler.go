package handler

import tb "gopkg.in/telebot.v3"

type CommandHandler interface {
	// Command string of bot Command
	Command() string
	// Description of Command
	Description() string
	// Handle function
	Handle(ctx tb.Context) error

	Middlewares() []tb.MiddlewareFunc
}
