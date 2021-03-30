package config

import (
	"fmt"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
	"text/template"
)

var ()

type RunType string

var (
	version = "master"
	commit  = "newest"
	date    = "today"

	ProjectName          string = "flowerss"
	BotToken             string
	Socks5               string
	TelegraphToken       []string
	TelegraphAccountName string
	TelegraphAuthorName  string = "rssbot"
	TelegraphAuthorURL   string

	// EnableTelegraph 是否启用telegraph
	EnableTelegraph       bool = false
	PreviewText           int  = 0
	DisableWebPagePreview bool = false
	Mysql                 MysqlConfig
	SQLitePath            string
	EnableMysql           bool = false

	// UpdateInterval rss抓取间隔
	UpdateInterval int = 10

	// ErrorThreshold rss源抓取错误阈值
	ErrorThreshold uint = 100

	// MessageTpl rss更新推送模版
	MessageTpl *template.Template

	// MessageMode telegram消息渲染模式
	MessageMode tb.ParseMode

	// TelegramEndpoint telegram bot 服务器地址，默认为空
	TelegramEndpoint string = tb.DefaultApiURL

	// UserAgent User-Agent
	UserAgent string

	// RunMode 运行模式 Release / Debug
	RunMode RunType = ReleaseMode

	// AllowUsers 允许使用bot的用户
	AllowUsers []int64

	// DBLogMode 是否打印数据库日志
	DBLogMode bool = false

	// 当 feed 有多个条目更新时，是否合并为一条消息
	MergeMessage bool = false

	// 使用列表发送消息时的新消息阈值
	MergeTolerance int = 2

	// 合并消息中单条消息模板，模式与 MessageMode 一致
	MessageItemTpl *template.Template
)

const (
	logo = `
   __ _                                
  / _| | _____      _____ _ __ ___ ___ 
 | |_| |/ _ \ \ /\ / / _ \ '__/ __/ __|
 |  _| | (_) \ V  V /  __/ |  \__ \__ \
 |_| |_|\___/ \_/\_/ \___|_|  |___/___/

`
	defaultMessageTplMode = tb.ModeHTML
	defaultMessageTpl     = `<b>{{.SourceTitle}}</b>{{ if .PreviewText }}
---------- Preview ----------
{{.PreviewText}}
-----------------------------
{{- end}}{{if .EnableTelegraph}}
<a href="{{.TelegraphURL}}">【预览】</a><a href="{{.RawLink}}">{{.ContentTitle}}</a>
{{- else }}
<a href="{{.RawLink}}">{{.ContentTitle}}</a>
{{- end }}
{{.Tags}}
`
	defaultMessageMarkdownTpl = `** {{.SourceTitle}} **{{ if .PreviewText }}
---------- Preview ----------
{{.PreviewText}}
-----------------------------
{{- end}}{{if .EnableTelegraph}}
[【预览】]({{.TelegraphURL}})[{{.ContentTitle}}]({{.RawLink}})
{{- else }}
[{{.ContentTitle}}]({{.RawLink}})
{{- end }}
{{.Tags}}
`
	//defaultMessageListItemTpl = `{{if .EnableTelegraph}}
	//[【预览】]({{.TelegraphURL}})[{{.ContentTitle}}]({{.RawLink}}) {{- else }}
	//[{{.ContentTitle}}]({{.RawLink}}){{- end }}`
	defaultMessageListItemTpl = `{{if .EnableTelegraph}}
<a href="{{.TelegraphURL}}">【预览】</a><a href="{{.RawLink}}">{{.ContentTitle}}</a>
{{- else }}
<a href="{{.RawLink}}">{{.ContentTitle}}</a>
{{- end }}`

	TestMode    RunType = "Test"
	ReleaseMode RunType = "Release"
)

// MysqlConfig mysql 配置
type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

type TplData struct {
	SourceTitle     string
	ContentTitle    string
	RawLink         string
	PreviewText     string
	TelegraphURL    string
	Tags            string
	EnableTelegraph bool
	IsItem          bool
}

func AppVersionInfo() (s string) {
	s = fmt.Sprintf("version %v, commit %v, built at %v", version, commit, date)
	return
}

// GetString get string config value by key
func GetString(key string) string {
	var value string
	if viper.IsSet(key) {
		value = viper.GetString(key)
	}

	return value
}
