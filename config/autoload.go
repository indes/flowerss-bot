package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

var (
	BotToken string
	Socks5   string
	Mysql    MysqlConfig
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
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     getInt(os.Getenv("MYSQL_PORT")),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DB:       os.Getenv("MYSQL_DB"),
	}
}

func getInt(s string) int {
	num, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	return num
}
