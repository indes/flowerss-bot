package session

// EntityType is a MessageEntity type.
type BotContextStoreKey string

const (
	MentionChat BotContextStoreKey = "mention_chat"
)

func (k BotContextStoreKey) String() string {
	return string(k)
}
