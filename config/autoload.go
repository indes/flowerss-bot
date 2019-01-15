package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

var (
	BotToken    string
	Socks5      string
	Mysql       MysqlConfig
	EnableMysql bool
)

type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

func init() {
	projectName := "rssflow"

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", projectName))              // call multiple times to add many search paths
	viper.AddConfigPath(fmt.Sprintf("/data/docker/config/%s", projectName)) // path to look for the config file in

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	BotToken = viper.GetString("token")
	Socks5 = viper.GetString("socks5")

	Mysql = MysqlConfig{
		Host:     viper.GetString("mysql.host"),
		Port:     getInt(viper.GetString("mysql.port")),
		User:     viper.GetString("mysql.user"),
		Password: viper.GetString("mysql.password"),
		DB:       viper.GetString("mysql.database"),
	}

	if Mysql.Host != "" {
		EnableMysql = true
	} else {
		EnableMysql = false
	}
}

func getInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func GetMysqlConnectingString() string {
	usr := viper.GetString("mysql.user")
	pwd := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	db := viper.GetString("mysql.database")
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=true", usr, pwd, host, db)
}
