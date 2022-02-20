package session

import "testing"

func TestBotContextStoreKey_String(t *testing.T) {
	tests := []struct {
		name string
		k    BotContextStoreKey
		want string
	}{
		{"mention", StoreKeyMentionChat, string(StoreKeyMentionChat)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
