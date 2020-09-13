package util

import (
	"github.com/indes/flowerss-bot/log"
	"net/http"
	"time"

	"github.com/indes/flowerss-bot/config"
	"golang.org/x/net/proxy"
)

var (
	HttpClient *http.Client
)

func clientInit() {
	httpTransport := &http.Transport{}
	HttpClient = &http.Client{Transport: httpTransport, Timeout: 15 * time.Second}
	// set proxy
	if config.Socks5 != "" {
		log.Infow("enable proxy",
			"socks5", config.Socks5,
		)

		dialer, err := proxy.SOCKS5("tcp", config.Socks5, nil, proxy.Direct)
		if err != nil {
			log.Fatal("Error creating dialer, aborting.")
		}
		httpTransport.Dial = dialer.Dial
	}
}
