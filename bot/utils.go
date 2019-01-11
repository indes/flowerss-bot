package bot

import "gopkg.in/tucnak/telebot.v2"

func SendError(c *telebot.Chat) {
	B.Send(c, "请输入正确的指令！")
}
