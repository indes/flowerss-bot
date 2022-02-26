package bot

import (
	"testing"

	"github.com/magiconair/properties/assert"
	tb "gopkg.in/telebot.v3"
)

// TestGetMentionFromMessage test GetMentionFromMessage
func TestGetMentionFromMessage(t *testing.T) {
	tests := []struct {
		name        string
		m           *tb.Message
		wantMention string
	}{
		{
			"null mention",
			&tb.Message{
				Text: "hello world",
			},
			"",
		},
		{
			"one text mention",
			&tb.Message{
				Text: "@hello world",
				Entities: []tb.MessageEntity{
					tb.MessageEntity{
						Type:   tb.EntityMention,
						Offset: 0,
						Length: 6,
					},
				},
			},
			"@hello",
		},
		{
			"multiple text mention", // get first mention
			&tb.Message{
				Text: "@hello @world!",
				Entities: []tb.MessageEntity{
					tb.MessageEntity{
						Type:   tb.EntityMention,
						Offset: 0,
						Length: 6,
					},
					tb.MessageEntity{
						Type:   tb.EntityMention,
						Offset: 7,
						Length: 7,
					},
				},
			},
			"@hello",
		},
		{
			"one caption mention",
			&tb.Message{
				Caption: "@hello world",
				CaptionEntities: []tb.MessageEntity{
					tb.MessageEntity{
						Type:   tb.EntityMention,
						Offset: 0,
						Length: 6,
					},
				},
			},
			"@hello",
		},
		{
			"multiple caption mention", // get first mention
			&tb.Message{
				Caption: "@hello @world!",
				CaptionEntities: []tb.MessageEntity{
					tb.MessageEntity{
						Type:   tb.EntityMention,
						Offset: 0,
						Length: 6,
					},
					tb.MessageEntity{
						Type:   tb.EntityMention,
						Offset: 7,
						Length: 7,
					},
				},
			},
			"@hello",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotMention := GetMentionFromMessage(tt.m)
				assert.Equal(t, gotMention, tt.wantMention)
			},
		)
	}
}
