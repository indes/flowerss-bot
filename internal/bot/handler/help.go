package handler

import (
	tb "gopkg.in/telebot.v3"
)

type Help struct {
}

func NewHelp() *Help {
	return &Help{}
}

func (h *Help) Command() string {
	return "/help"
}

func (h *Help) Description() string {
	return "帮助"
}

func (h *Help) Handle(ctx tb.Context) error {
	message := `
	命令：
	/sub 订阅源
	/unsub  取消订阅
	/list 查看当前订阅源
	/set 设置订阅
	/check 检查当前订阅
	/setfeedtag 设置订阅标签
	/setinterval 设置订阅刷新频率
	/activeall 开启所有订阅
	/pauseall 暂停所有订阅
	/help 帮助
	/import 导入 OPML 文件
	/export 导出 OPML 文件
	/unsuball 取消所有订阅
	详细使用方法请看：https://github.com/indes/flowerss-bot
	`
	return ctx.Send(message)
}

func (h *Help) Middlewares() []tb.MiddlewareFunc {
	return nil
}
