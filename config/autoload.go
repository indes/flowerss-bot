package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strconv"
)

var (
	BotToken        string
	Socks5          string
	TelegraphToken  []string
	EnableTelegraph bool
	Mysql           MysqlConfig
	SQLitePath      string
	EnableMysql     bool
	UpdateInterval  int  = 10
	ErrorThreshold  uint = 100
)

type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

func init() {

	logo := `
   __ _                                
  / _| | _____      _____ _ __ ___ ___ 
 | |_| |/ _ \ \ /\ / / _ \ '__/ __/ __|
 |  _| | (_) \ V  V /  __/ |  \__ \__ \
 |_| |_|\___/ \_/\_/ \___|_|  |___/___/

`
	fmt.Println(logo)
	telegramTokenCli := flag.String("b", "", "Telegram Bot Token")
	telegraphTokenCli := flag.String("t", "", "Telegraph API Token")
	dbPathCli := flag.String("dbpath", "", "Telegraph API Token")
	errorThresholdCli := flag.Int("threshold", 0, "Error Threshold")
	socks5Cli := flag.String("s", "", "Socks5 Proxy")
	intervalCli := flag.Int("i", 0, "Update Interval")
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
