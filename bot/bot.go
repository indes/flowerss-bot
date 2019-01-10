package bot

import (
	"github.com/indes/go-rssbot/bot/env"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"time"
)

var (
	botToken    = env.BotToken
	socks5Proxy = env.Socks5
)

func init() {
}

func Start() {
	log.Printf("Token: %s Proxy: %s\n", botToken, socks5Proxy)

	// creat bot
	dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
	if err != nil {
		log.Fatal("Error creating dialer, aborting.")
	}

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial

	b, err := tb.NewBot(tb.Settings{
		Token: botToken,
		// You can also set custom API URL. If field is empty it equals to "https://api.telegram.org"
		//URL:    "http://195.129.111.17:8012",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: httpClient,
	})

	err = makeHandle(b)
	b.Start()
}

func makeHandle(b *tb.Bot) error {

	b.Handle("/hello", func(m *tb.Message) {
		log.Println(m.Text)
		_, _ = b.Send(m.Sender, "hello world!")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {

	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		// photos only
	})

	b.Handle(tb.OnChannelPost, func(m *tb.Message) {
		// channel posts only
	})

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		// incoming inline queries
	})

	return nil
}
