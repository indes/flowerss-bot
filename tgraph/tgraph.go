package tgraph

import (
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
	client     *telegraph.Client
)

func init() {
	if config.EnableTelegraph {
		log.Println("Telegraph Enable")
		telegraph.Verbose = verbose
		var err error
		client, err = telegraph.Load(authToken)
		if err != nil {
			log.Fatal("* Load error: %s", err)
		}
	}

}
