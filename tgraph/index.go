package tgraph

import (
	"github.com/meinside/telegraph-go"
	"log"
)

const (
	//verbose = false

	html = `<h1>hello</h1>`
)

var (
	client *telegraph.Client
)

func init() {
	telegraph.Verbose = verbose
	var err error
	client, err = telegraph.Load(authToken)
	if err != nil {
		log.Fatal("* Load error: %s", err)
	}
}
