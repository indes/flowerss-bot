package handler

import tb "gopkg.in/telebot.v3"

type Import struct {
}

func NewImport() *Import {
	return &Import{}
}

func (i *Import) Command() string {
	return "/import"
}

func (i *Import) Description() string {
	return "导入OPML文件"
}

func (i *Import) Handle(ctx tb.Context) error {
	reply := "请直接发送OPML文件，如果需要为频道导入OPML，请在发送文件的时候附上channel id，例如@telegram"
	return ctx.Reply(reply)
}

func (i *Import) Middlewares() []tb.MiddlewareFunc {
	return nil
}
