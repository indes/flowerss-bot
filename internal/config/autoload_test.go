package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
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
		//{
		//	"markdown",
		//	fields{SourceTitle: "[aaa](qq) *123*"},
		//	args{telebot.ModeMarkdown},
		//	"** \\[aaa](qq) \\*123\\* **\n[]()",
		//	false,
		//},
		{"HTML Mode",
			fields{SourceTitle: "[aaa] *123*", ContentTitle: "google", RawLink: "https://google.com"},
			args{telebot.ModeHTML},
			"<b>[aaa] *123*</b>\n<a href=\"https://google.com\">google</a>",
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

			assert.Equal(t1, tt.want, got)
			assert.Equal(t1, err != nil, tt.wantErr)
		})
	}
}

func TestTplData_replaceHTMLTags(t1 *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"case1", "<hello>", "&lt;hello&gt;"},
		{"case2", "<\"hello\">", "&lt;&quot;hello&quot;&gt;"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TplData{}

			got := t.replaceHTMLTags(tt.arg)
			assert.Equal(t1, tt.want, got)

		})
	}
}
