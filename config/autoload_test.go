package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
	"testing"
)

func TestTplData_Render(t1 *testing.T) {
	type fields struct {
		SourceTitle     string
		ContentTitle    string
		RawLink         string
		PreviewText     string
		TelegraphURL    string
		Tags            string
		EnableTelegraph bool
	}
	type args struct {
		mode telebot.ParseMode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"markdown",
			fields{SourceTitle: "[aaa](qq) *123*"},
			args{telebot.ModeMarkdown},
			"** \\[aaa](qq) \\*123\\* **\n[]()",
			false,
		},
		{"HTML Mode",
			fields{SourceTitle: "[aaa] *123*"},
			args{telebot.ModeHTML},
			"** [aaa] *123* **\n[]()",
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TplData{
				SourceTitle:     tt.fields.SourceTitle,
				ContentTitle:    tt.fields.ContentTitle,
				RawLink:         tt.fields.RawLink,
				PreviewText:     tt.fields.PreviewText,
				TelegraphURL:    tt.fields.TelegraphURL,
				Tags:            tt.fields.Tags,
				EnableTelegraph: tt.fields.EnableTelegraph,
			}
			got, err := t.Render(tt.args.mode)

			assert.Equal(t1, got, tt.want)
			assert.Equal(t1, err != nil, tt.wantErr)
		})
	}
}
