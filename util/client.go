package util

import (
	"github.com/indes/flowerss-bot/config"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
)

var (
	HttpClient *http.Client
)

func clientInit() {
	httpTransport := &http.Transport{}
	HttpClient = &http.Client{Transport: httpTransport}
	// set proxy
	if config.Socks5 != "" {
		log.Printf("Proxy: %s\n", config.Socks5)

		dialer, err := proxy.SOCKS5("tcp", config.Socks5, nil, proxy.Direct)
		if err != nil {
			log.Fatal("Error creating dialer, aborting.")
		}
		httpTransport.Dial = dialer.Dial
	}
}
