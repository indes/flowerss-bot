package bot

import (
	"github.com/indes/go-rssbot/config"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	botToken            = config.BotToken
	socks5Proxy         = config.Socks5
	B           *tb.Bot = nil
)

func init() {
	log.Printf("Token: %s Proxy: %s\n", botToken, socks5Proxy)

	// creat bot
	dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
	if err != nil {
		log.Fatal("Error creating dialer, aborting.")
	}

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial

	B, err = tb.NewBot(tb.Settings{
		Token: botToken,
		// You can also set custom API URL. If field is empty it equals to "https://api.telegram.org"
		// URL:    "http://195.129.111.17:8012",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: httpClient,
	})

	if err != nil {
		log.Fatal(err)
		return
	}
}

func Start() {

	register(B)
	makeHandle(B)

	B.Start()
}

func register(b *tb.Bot) {

}

func makeHandle(b *tb.Bot) {

	b.Handle("/hello", func(m *tb.Message) {
		log.Println(m.Text)
		_, _ = b.Send(m.Sender, "hello world!")
	})

	b.Handle("/sub", func(m *tb.Message) {
		//log.Fatal(m.Text)
		msg := strings.Split(m.Text, " ")

		if len(msg) > 2 {
			SendError(m.Chat)
		} else {
			url := msg[1]
			if CheckUrl(url) {
				registFeed(m.Chat, url);
			} else {
				SendError(m.Chat)
			}
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {

	})
}
