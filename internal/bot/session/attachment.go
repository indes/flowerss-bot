package session

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Attachment 透传信息，用于bot按钮信息透传
type Attachment struct {
	UserID   int64 `json:"user_id"`
	SourceID uint  `json:"source_id"`
}

// Marshal 编码成json字符串
func (a *Attachment) Marshal() string {
	str, _ := json.MarshalToString(a)
	return str
}

// ParseAttachment 从json字符串解析透传信息
func UnmarshalAttachment(data string) (*Attachment, error) {
	a := &Attachment{}
	if err := json.UnmarshalFromString(data, a); err != nil {
		return nil, err
	}
	return a, nil
}
