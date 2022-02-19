package bot

import tb "gopkg.in/telebot.v3"

type newContentMessage struct {
}

func (c *newContentMessage) Send() (*tb.Message, error) {
	return nil, nil
}
