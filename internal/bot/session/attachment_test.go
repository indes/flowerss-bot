package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttachment(t *testing.T) {
	t.Run(
		"encode and decode", func(t *testing.T) {
			a := &Attachment{UserID: 123, SourceID: 321}
			data := a.Marshal()
			a2, err := UnmarshalAttachment(data)
			assert.Nil(t, err)
			assert.Equal(t, a, a2)
		},
	)
}
