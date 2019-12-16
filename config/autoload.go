package config

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
)

var (
	BotToken        string
	Socks5          string
	TelegraphToken  []string
	EnableTelegraph bool
	PreviewText     int = 0
	Mysql           MysqlConfig
	SQLitePath      string
	EnableMysql     bool
	UpdateInterval  int  = 10
	ErrorThreshold  uint = 100
	MessageTpl      *template.Template
	MessageMode     tb.ParseMode
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
{{- end }}`
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
	EnableTelegraph bool
}

func validateTPL() {
	testData := []TplData{
		TplData{
			"RSS 源标识 - 无预览无telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"",
			"",
			false,
		},
		TplData{
			"RSS源标识 - 有预览无telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"这里是很长很长很长的消息预览字数补丁紫薯补丁紫薯补丁紫薯补丁紫薯补丁",
			"",
			false,
		},
		TplData{
			"RSS源标识 - 有预览有telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"这里是很长很长很长的消息预览字数补丁紫薯补丁紫薯补丁紫薯补丁紫薯补丁",
			"https://telegra.ph/markdown-07-07",
			true,
		},
	}

	var buf []byte
	w := bytes.NewBuffer(buf)

	for _, d := range testData {
		w.Reset()
		fmt.Println("\n////////////////////////////////////////////")
		if err := MessageTpl.Execute(os.Stdout, d); err != nil {
			log.Fatal(err)
		}
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

	if viper.IsSet("message_tpl_mode") {
		switch strings.ToLower(viper.GetString("message_tpl_mode")) {
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

func init() {

	telegramTokenCli := flag.String("b", "", "Telegram Bot Token")
	telegraphTokenCli := flag.String("t", "", "Telegraph API Token")
	previewTextCli := flag.Int("p", 0, "Preview Text Length")
	dbPathCli := flag.String("dbpath", "", "SQLite DB Path")
	errorThresholdCli := flag.Int("threshold", 0, "Error Threshold")
	socks5Cli := flag.String("s", "", "Socks5 Proxy")
	intervalCli := flag.Int("i", 0, "Update Interval")
	testTpl := flag.Bool("testtpl", false, "Test Template")
	flag.Parse()

	projectName := "flowerss-bot"

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", projectName))              // call multiple times to add many search paths
	viper.AddConfigPath(fmt.Sprintf("/data/docker/config/%s", projectName)) // path to look for the config file in

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

	if *telegramTokenCli == "" {
		BotToken = viper.GetString("bot_token")

	} else {
		BotToken = *telegramTokenCli
	}

	if *socks5Cli == "" {
		Socks5 = viper.GetString("socks5")
	} else {
		Socks5 = *socks5Cli
	}

	if *telegraphTokenCli == "" {
		if viper.IsSet("telegraph_token") {
			EnableTelegraph = true

			TelegraphToken = viper.GetStringSlice("telegraph_token")
		} else {
			EnableTelegraph = false
		}
	} else {
		EnableTelegraph = true
		TelegraphToken = append(TelegraphToken, *telegraphTokenCli)
	}

	if *previewTextCli == 0 {
		if viper.IsSet("preview_text") {
			PreviewText = viper.GetInt("preview_text")
		} else {
			PreviewText = 0
		}
	} else {
		PreviewText = 0
	}

	if *errorThresholdCli == 0 {
		if viper.IsSet("error_threshold") {
			ErrorThreshold = uint(viper.GetInt("error_threshold"))
		} else {
			ErrorThreshold = 100
		}
	} else {
		ErrorThreshold = uint(*errorThresholdCli)
	}

	if *intervalCli == 0 {
		if viper.IsSet("update_interval") {
			UpdateInterval = viper.GetInt("update_interval")
		} else {
			UpdateInterval = 10
		}
	} else {
		UpdateInterval = *intervalCli
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
		if *dbPathCli == "" {
			if viper.IsSet("sqlite.path") {
				SQLitePath = viper.GetString("sqlite.path")
			} else {
				SQLitePath = "data.db"
			}
		} else {
			SQLitePath = *dbPathCli
		}
		log.Println("DB Path: ", SQLitePath)
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
