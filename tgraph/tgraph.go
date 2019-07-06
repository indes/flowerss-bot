package tgraph

import (
	"fmt"
	"github.com/indes/flowerss-bot/config"
	"github.com/meinside/telegraph-go"
	"log"
)

const (
	//verbose = false
	html = `<h1>hello</h1>`
)

var (
	authToken  = config.TelegraphToken
	authorUrl  = "https://github.com/indes/flowerss-bot"
	authorName = "flowerss"
	verbose    = false
	//client     *telegraph.Client
	clientPool []*telegraph.Client
)

func init() {
	if config.EnableTelegraph {
		log.Println("Telegraph Enabled, Token len: ", len(authToken), "Token: ", authToken)

		telegraph.Verbose = verbose

		for _, t := range authToken {
			client, err := telegraph.Load(t)
			if err != nil {
				log.Println(fmt.Sprintf("Telegraph load error: %s token: %s", err, t))
			} else {
				clientPool = append(clientPool, client)
			}
		}
		if len(clientPool) == 0 {
			config.EnableTelegraph = false
			log.Println("Telegraph token error, telegraph disabled")

		}

	}

}
