package session

import tb "gopkg.in/telebot.v3"

// EntityType is a MessageEntity type.
type BotContextStoreKey string

const (
	StoreKeyMentionChat BotContextStoreKey = "mention_chat"
)

func (k BotContextStoreKey) String() string {
	return string(k)
}

func GetMentionChatFromCtxStore(ctx tb.Context) (*tb.Chat, bool) {
	v := ctx.Get(StoreKeyMentionChat.String())
	if v == nil {
		return nil, false
	}

	mentionChat, ok := v.(*tb.Chat)
	if !ok {
		return nil, false
	}
	return mentionChat, true
}
