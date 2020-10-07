package bot

import (
	"github.com/indes/flowerss-bot/config"
	"github.com/magiconair/properties/assert"
	tb "gopkg.in/tucnak/telebot.v2"
	"testing"
)

// TestGetMentionFromMessage test GetMentionFromMessage
func TestGetMentionFromMessage(t *testing.T) {
	tests := []struct {
		name        string
		m           *tb.Message
		wantMention string
	}{
		{"null mention",
			&tb.Message{
				Text: "hello world",
			},
			"",
		},
		{"one text mention",
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
		{"multiple text mention", // get first mention
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
		{"one caption mention",
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
		{"multiple caption mention", // get first mention
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
		t.Run(tt.name, func(t *testing.T) {
			gotMention := GetMentionFromMessage(tt.m)
			assert.Equal(t, gotMention, tt.wantMention)
		})
	}
}

// Test_isUserAllowed test isUserAllowed
func Test_isUserAllowed(t *testing.T) {
	tests := []struct {
		name string
		upd  *tb.Update
		want bool
	}{
		{
			"错误消息1",
			nil,
			false,
		},
		{
			"错误消息2",
			&tb.Update{},
			false,
		},
		{
			"anybody",
			&tb.Update{
				Message: &tb.Message{
					Sender: &tb.User{
						ID: 123,
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUserAllowed(tt.upd)
			assert.Equal(t, got, tt.want)
		})
	}

	config.AllowUsers = append(config.AllowUsers, 123)

	tests = []struct {
		name string
		upd  *tb.Update
		want bool
	}{
		{
			"白名单用户",
			&tb.Update{
				Message: &tb.Message{
					Sender: &tb.User{
						ID: 123,
					},
				},
			},
			true,
		},
		{
			"非白名单用户",
			&tb.Update{
				Message: &tb.Message{
					Sender: &tb.User{
						ID: 321,
					},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUserAllowed(tt.upd)
			assert.Equal(t, got, tt.want)
		})
	}
}
