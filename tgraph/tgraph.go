package tgraph

import (
	"fmt"
	"log"

	"github.com/indes/flowerss-bot/config"
	"github.com/indes/telegraph-go"
)

const (
	//verbose = false
	htmlContent = `<h1>hello</h1>`
)

var (
	authToken   = config.TelegraphToken
	socks5Proxy = config.Socks5
	authorUrl   = "https://github.com/indes/flowerss-bot"
	authorName  = "flowerss"
	verbose     = false
	//client     *telegraph.Client
	clientPool []*telegraph.Client
)

func init() {
	if config.EnableTelegraph {
		log.Println("Telegraph Enabled, Token len: ", len(authToken), "Token: ", authToken)

		telegraph.Verbose = verbose

		for _, t := range authToken {
			client, err := telegraph.Load(t, socks5Proxy)
			if err != nil {
				log.Println(fmt.Sprintf("Telegraph load error: %s token: %s", err, t))
			} else {
				clientPool = append(clientPool, client)
			}
		}

		if len(clientPool) == 0 {
			if config.TelegraphAccountName == "" {
				config.EnableTelegraph = false
				log.Println("Telegraph token error, telegraph disabled")
			} else if len(authToken) == 0 {
				// create account
				if client, err := telegraph.Create(config.TelegraphAccountName, config.TelegraphAuthorName, config.TelegraphAuthorURL, config.Socks5); err != nil {
					log.Println("create telegraph account fail: ", err)
				} else {
					clientPool = append(clientPool, client)
				}
			}
		}
	}

}
