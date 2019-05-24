package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

var (
	BotToken        string
	Socks5          string
	TelegraphToken  string
	EnableTelegraph bool
	Mysql           MysqlConfig
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
	projectName := "flowerss-bot"

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", projectName))              // call multiple times to add many search paths
	viper.AddConfigPath(fmt.Sprintf("/data/docker/config/%s", projectName)) // path to look for the config file in

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	BotToken = viper.GetString("bot_token")
	Socks5 = viper.GetString("socks5")

	if viper.IsSet("telegraph_token") {
		EnableTelegraph = true
		TelegraphToken = viper.GetString("telegraph_token")
	} else {
		EnableTelegraph = false
	}

	if viper.IsSet("update_interval") {
		UpdateInterval = viper.GetInt("update_interval")
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
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", usr, pwd, host, port, db)
}
