package message

import tb "gopkg.in/telebot.v3"

// MentionFromMessage get message mention
func MentionFromMessage(m *tb.Message) (mention string) {
	if m.Text != "" {
		for _, entity := range m.Entities {
			if entity.Type != tb.EntityMention {
				continue
			}
			return m.Text[entity.Offset : entity.Offset+entity.Length]
		}
	}

	for _, entity := range m.CaptionEntities {
		if entity.Type == tb.EntityMention {
			continue
		}
		return m.Caption[entity.Offset : entity.Offset+entity.Length]
	}
	return ""
}
