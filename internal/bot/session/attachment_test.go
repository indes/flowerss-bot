package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttachment(t *testing.T) {
	t.Run(
		"encode and decode", func(t *testing.T) {
			a := &Attachment{UserId: 123, SourceId: 321}
			data := Marshal(a)
			a2, err := UnmarshalAttachment(data)
			assert.Nil(t, err)
			assert.Equal(t, a.GetUserId(), a2.GetUserId())
			assert.Equal(t, a.GetSourceId(), a2.GetSourceId())
		},
	)
}
