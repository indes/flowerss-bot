package tgraph

import (
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/log"

	"github.com/indes/telegraph-go"
)

var (
	authToken   = config.TelegraphToken
	socks5Proxy = config.Socks5
	clientPool  []*telegraph.Client
)

func init() {
	if config.EnableTelegraph {
		log.Infof("telegraph enabled, count %d, %#v", len(authToken), authToken)
		for _, t := range authToken {
			client, err := telegraph.Load(t, socks5Proxy)
			if err != nil {
				log.Errorf("telegraph load %s failed, %v", t, err)
			} else {
				clientPool = append(clientPool, client)
			}
		}

		if len(clientPool) == 0 {
			if config.TelegraphAccountName == "" {
				config.EnableTelegraph = false
				log.Error("telegraph token error, telegraph disabled")
			} else if len(authToken) == 0 {
				// create account
				client, err := telegraph.Create(
					config.TelegraphAccountName,
					config.TelegraphAuthorName,
					config.TelegraphAuthorURL,
					config.Socks5,
				)

				if err != nil {
					config.EnableTelegraph = false
					log.Errorf("create telegraph account failed, %v", err)
				}

				clientPool = append(clientPool, client)
				log.Infof("create telegraph account success, token %v", client.AccessToken)
			}
		}
	}
}
