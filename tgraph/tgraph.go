package tgraph

import (
	"github.com/indes/rssflow/config"
	"github.com/meinside/telegraph-go"
	"log"
)

const (
	//verbose = false
	html = `<h1>hello</h1>`
)

var (
	authToken  = config.TelegraphToken
	authorUrl  = "https://github.com/indes/rssflow"
	authorName = "RSSFlow"
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
