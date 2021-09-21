package bot

import "gopkg.in/tucnak/telebot.v2"

type newContentMessage struct {
}

func (c *newContentMessage) Send() (*telebot.Message, error) {
	return nil, nil
}
