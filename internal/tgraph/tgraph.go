package tgraph

import (
	"github.com/indes/flowerss-bot/internal/config"

	"github.com/indes/telegraph-go"
	"go.uber.org/zap"
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
		zap.S().Infow("telegraph enabled",
			"token count", len(authToken),
			"token list", authToken,
		)

		telegraph.Verbose = verbose

		for _, t := range authToken {
			client, err := telegraph.Load(t, socks5Proxy)
			if err != nil {
				zap.S().Errorw("telegraph load error",
					"error", err,
					"token", t,
				)
			} else {
				clientPool = append(clientPool, client)
			}
		}

		if len(clientPool) == 0 {
			if config.TelegraphAccountName == "" {
				config.EnableTelegraph = false
				zap.S().Error("telegraph token error, telegraph disabled")
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
					zap.S().Errorw("create telegraph account fail, telegraph disabled", "error", err)
				}

				clientPool = append(clientPool, client)
				zap.S().Infow("create telegraph account success",
					"telegraph token", client.AccessToken)

			}
		}
	}
}
