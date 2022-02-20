package message

import (
	"regexp"

	tb "gopkg.in/telebot.v3"
)

// MentionFromMessage get message mention
func MentionFromMessage(m *tb.Message) string {
	if m.Text != "" {
		for _, entity := range m.Entities {
			if entity.Type != tb.EntityMention {
				continue
			}
			return m.Text[entity.Offset : entity.Offset+entity.Length]
		}
	}

	for _, entity := range m.CaptionEntities {
		if entity.Type != tb.EntityMention {
			continue
		}
		return m.Caption[entity.Offset : entity.Offset+entity.Length]
	}
	return ""
}

var relaxUrlMatcher = regexp.MustCompile(`^(https?://.*?)($| )`)

// URLFromMessage get message url
func URLFromMessage(m *tb.Message) string {
	for _, entity := range m.Entities {
		if entity.Type == tb.EntityURL {
			return m.Text[entity.Offset : entity.Offset+entity.Length]
		}
	}

	var payloadMatching = relaxUrlMatcher.FindStringSubmatch(m.Payload)
	if len(payloadMatching) > 0 {
		return payloadMatching[0]
	}
	return ""
}
