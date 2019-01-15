package bot

import (
	"fmt"
	"github.com/indes/rssflow/config"
	"github.com/indes/rssflow/model"
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

	if socks5Proxy != "" {
		dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Fatal("Error creating dialer, aborting.")
		}

		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial

		// creat bot
		B, err = tb.NewBot(tb.Settings{
			Token:  botToken,
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
			Client: httpClient,
		})

		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		var err error
		// creat bot
		B, err = tb.NewBot(tb.Settings{
			Token:  botToken,
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			log.Fatal(err)
			return
		}
	}

}

func Start() {
	makeHandle()
	B.Start()
}

func makeHandle() {

	B.Handle("/start", func(m *tb.Message) {
		user := model.FindOrInitUser(m.Chat.ID)
		_, _ = B.Send(m.Sender, fmt.Sprintf("hello, %d", user.ID))
	})

	B.Handle("/sub", func(m *tb.Message) {

		msg := strings.Split(m.Text, " ")

		if len(msg) != 2 {
			SendError(m.Chat)
		} else {
			url := msg[1]
			if CheckUrl(url) {
				registFeed(m.Chat, url)
			} else {
				SendError(m.Chat)
			}
		}
	})

	B.Handle("/list", func(m *tb.Message) {
		sources, _ := model.GetSourcesByUserID(m.Sender.ID)
		
		message := "目前的订阅源：\n"
		for index, source := range sources {
			message = message + fmt.Sprintf("[[%d]] [%s](%s)\n", index+1, source.Title, source.Link)
		}
		_, _ = B.Send(m.Sender, message, &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})

	})

	B.Handle("/unsub", func(m *tb.Message) {
		msg := strings.Split(m.Text, " ")

		if len(msg) != 2 {
			SendError(m.Chat)
		} else {
			url := msg[1]
			if CheckUrl(url) {
				source, _ := model.GetSourceByUrl(url)
				if source == nil {
					_, _ = B.Send(m.Sender, "未订阅该RSS源")
				} else {
					err := model.UnsubByUserIDAndSource(m.Sender.ID, source)
					if err == nil {
						_, _ = B.Send(m.Sender, "退定成功！")
					} else {
						_, err = B.Send(m.Sender, err.Error())
					}
				}
			} else {
				SendError(m.Chat)
			}
		}

	})

	B.Handle("/ping", func(m *tb.Message) {

		_, _ = B.Send(m.Sender, "pong")
	})

	B.Handle("/test", func(m *tb.Message) {

		message := `
*bold text*
_italic text_
[inline URL](http://www.example.com/)
[inline mention of a user](tg://user?id=123456789)
`

		_, err := B.Send(m.Sender, message, &tb.SendOptions{
			ParseMode: tb.ModeMarkdown,
		})
		log.Println(err)
	})
	B.Handle(tb.OnText, func(m *tb.Message) {

	})
}
