package config

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	ProjectName           string = "flowerss"
	BotToken              string
	Socks5                string
	TelegraphToken        []string
	EnableTelegraph       bool
	PreviewText           int = 0
	DisableWebPagePreview bool
	Mysql                 MysqlConfig
	SQLitePath            string
	EnableMysql           bool
	UpdateInterval        int  = 10
	ErrorThreshold        uint = 100
	MessageTpl            *template.Template
	MessageMode           tb.ParseMode
	TelegramEndpoint      string
	UserAgent             string
)

const (
	logo = `
   __ _                                
  / _| | _____      _____ _ __ ___ ___ 
 | |_| |/ _ \ \ /\ / / _ \ '__/ __/ __|
 |  _| | (_) \ V  V /  __/ |  \__ \__ \
 |_| |_|\___/ \_/\_/ \___|_|  |___/___/

`
	defaultMessageTplMode = "md"
	defaultMessageTpl     = `** {{.SourceTitle}} **{{ if .PreviewText }}
---------- Preview ----------
{{.PreviewText}}
-----------------------------
{{- end}}{{if .EnableTelegraph}}
{{.ContentTitle}} [Telegraph]({{.TelegraphURL}}) | [原文]({{.RawLink}})
{{- else }}
[{{.ContentTitle}}]({{.RawLink}})
{{- end }}
{{.Tags}}
`
)

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
}

func init() {

	workDirFlag := flag.String("d", "./", "work directory of flowerss")
	configFile := flag.String("c", "", "config file of flowerss")
	printVersionFlag := flag.Bool("v", false, "prints flowerss-bot version")

	testTpl := flag.Bool("testtpl", false, "test template")

	//telegramTokenCli := flag.String("b", "", "Telegram Bot Token")
	//telegraphTokenCli := flag.String("t", "", "Telegraph API Token")
	//previewTextCli := flag.Int("p", 0, "Preview Text Length")
	//DisableWebPagePreviewCli := flag.Bool("disable_web_page_preview", false, "Disable Web Page Preview")
	//dbPathCli := flag.String("dbpath", "", "SQLite DB Path")
	//errorThresholdCli := flag.Int("threshold", 0, "Error Threshold")
	//socks5Cli := flag.String("s", "", "Socks5 Proxy")
	//intervalCli := flag.Int("i", 0, "Update Interval")
	//TelegramEndpointCli := flag.String("endpoint", "", "Custom Telegram Endpoint")

	flag.Parse()

	if *printVersionFlag {
		// print version
		fmt.Printf("version %v, commit %v, built at %v", version, commit, date)
		os.Exit(0)
	}

	workDir := filepath.Clean(*workDirFlag)

	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigFile(filepath.Join(workDir, "config.yml"))
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	initTPL()
	if *testTpl {
		validateTPL()
		os.Exit(0)
	}

	fmt.Println(logo)

	BotToken = viper.GetString("bot_token")
	Socks5 = viper.GetString("socks5")
	UserAgent = viper.GetString("user_agent")

	if viper.IsSet("telegraph_token") {
		EnableTelegraph = true
		TelegraphToken = viper.GetStringSlice("telegraph_token")
	} else {
		EnableTelegraph = false
	}

	if viper.IsSet("preview_text") {
		PreviewText = viper.GetInt("preview_text")
	} else {
		PreviewText = 0
	}

	DisableWebPagePreview = viper.GetBool("disable_web_page_preview")

	if viper.IsSet("telegram.endpoint") {
		TelegramEndpoint = viper.GetString("telegram.endpoint")
	} else {
		TelegramEndpoint = tb.DefaultApiURL
	}

	if viper.IsSet("error_threshold") {
		ErrorThreshold = uint(viper.GetInt("error_threshold"))
	} else {
		ErrorThreshold = 100
	}

	if viper.IsSet("update_interval") {
		UpdateInterval = viper.GetInt("update_interval")
	} else {
		UpdateInterval = 10
	}

	if viper.IsSet("mysql.host") {
		EnableMysql = true
		Mysql = MysqlConfig{
			Host:     viper.GetString("mysql.host"),
			Port:     viper.GetInt("mysql.port"),
			User:     viper.GetString("mysql.user"),
			Password: viper.GetString("mysql.password"),
			DB:       viper.GetString("mysql.database"),
		}
	} else {
		EnableMysql = false
	}

	if !EnableMysql {
		if viper.IsSet("sqlite.path") {
			SQLitePath = viper.GetString("sqlite.path")
		} else {
			SQLitePath = filepath.Join(workDir, "data.db")
		}
		log.Println("DB Path:", SQLitePath)
		// 判断并创建SQLite目录
		dir := path.Dir(SQLitePath)
		_, err := os.Stat(dir)
		if err != nil {
			err := os.MkdirAll(dir, os.ModeDir)
			if err != nil {
				log.Printf("mkdir failed![%v]\n", err)
			}
		}
	}

}

func (t TplData) Render(mode tb.ParseMode) (string, error) {

	var buf []byte
	wb := bytes.NewBuffer(buf)

	if mode == tb.ModeMarkdown {
		mkd := regexp.MustCompile("[\\*\\[\\]`_]")
		t.SourceTitle = mkd.ReplaceAllString(t.SourceTitle, " ")
		t.ContentTitle = mkd.ReplaceAllString(t.ContentTitle, " ")
		t.PreviewText = mkd.ReplaceAllString(t.PreviewText, " ")
	}

	if err := MessageTpl.Execute(wb, t); err != nil {
		return "", err
	}

	return strings.TrimSpace(string(wb.Bytes())), nil
}

func validateTPL() {
	testData := []TplData{
		TplData{
			"RSS 源标识 - 无预览无telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"",
			"",
			"",
			false,
		},
		TplData{
			"RSS源标识 - 有预览无telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"这里是很长很长很长的消息预览字数补丁紫薯补丁紫薯补丁紫薯补丁紫薯补丁[1](123)",
			"",
			"#标签",
			false,
		},
		TplData{
			"RSS源标识 - 有预览有telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"这里是很长很长很长的消息预览字数补丁紫薯补丁紫薯补丁紫薯补丁紫薯补丁",
			"https://telegra.ph/markdown-07-07",
			"#标签1 #标签2",
			true,
		},
	}

	for _, d := range testData {
		fmt.Println("\n////////////////////////////////////////////")
		fmt.Println(d.Render(MessageMode))
	}
	fmt.Println("\n////////////////////////////////////////////")
}

func initTPL() {

	var tplMsg string
	if viper.IsSet("message_tpl") {
		tplMsg = viper.GetString("message_tpl")
	} else {
		tplMsg = defaultMessageTpl
	}
	MessageTpl = template.Must(template.New("message").Parse(tplMsg))

	if viper.IsSet("message_mode") {
		switch strings.ToLower(viper.GetString("message_mode")) {
		case "md", "markdown":
			MessageMode = tb.ModeMarkdown
		case "html":
			MessageMode = tb.ModeHTML
		default:
			MessageMode = tb.ModeDefault
		}
	} else {
		MessageMode = tb.ModeMarkdown
	}

}

func getInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func (m *MysqlConfig) GetMysqlConnectingString() string {
	usr := m.User
	pwd := m.Password
	host := m.Host
	port := m.Port
	db := m.DB
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", usr, pwd, host, port, db)
}
