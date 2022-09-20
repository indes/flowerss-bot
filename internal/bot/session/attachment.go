package session

import (
	"encoding/hex"

	"google.golang.org/protobuf/proto"

	"github.com/indes/flowerss-bot/internal/log"
)

// Marshal 编码成字符串
func Marshal(a *Attachment) string {
	bytes, err := proto.Marshal(a)
	if err != nil {
		log.Errorf("marshal attachment failed, %v", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

// UnmarshalAttachment 从字符串解析透传信息
func UnmarshalAttachment(data string) (*Attachment, error) {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	a := &Attachment{}
	if err := proto.Unmarshal(bytes, a); err != nil {
		return nil, err
	}
	return a, nil
}
