package bot

import (
	"bytes"
	"fmt"
	"github.com/indes/rssflow/bot/fsm"
	"github.com/indes/rssflow/config"
	"github.com/indes/rssflow/model"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	botToken                           = config.BotToken
	socks5Proxy                        = config.Socks5
	UserState   map[int]fsm.UserStatus = make(map[int]fsm.UserStatus)
	//B bot
	B *tb.Bot
)

func init() {

	if socks5Proxy != "" {
		log.Printf("Bot Token: %s Proxy: %s\n", botToken, socks5Proxy)

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
		log.Printf("Bot Token: %s", botToken)

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

//Start run bot
func Start() {
	makeHandle()
	B.Start()
}

func makeHandle() {
	toggleNoticeKey := tb.InlineButton{
		Unique: "toggle_notice",
		Text:   "切换通知",
	}
	toggleTelegraphKey := tb.InlineButton{
		Unique: "toggle_telegraph",
		Text:   "切换 Telegraph",
	}
	feedSettingKeys := [][]tb.InlineButton{}
	B.Handle(&toggleNoticeKey, func(c *tb.Callback) {
		// on inline button pressed (callback!)

		// always respond!
		B.Respond(c, &tb.CallbackResponse{})
	})
	B.Handle("/key", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		toggleNoticeKey.Text = "开启通知"

		feedSettingKeys = append(
			feedSettingKeys,
			[]tb.InlineButton{toggleNoticeKey},
		)
		_, _ = B.Send(m.Sender, "Hello!", &tb.ReplyMarkup{
			InlineKeyboard: feedSettingKeys,
		})
	})

	B.Handle("/start", func(m *tb.Message) {
		user := model.FindOrInitUser(m.Chat.ID)
		log.Printf("/start %d", user.ID)
		_, _ = B.Send(m.Sender, fmt.Sprintf("hello"))
	})

	B.Handle("/export", func(m *tb.Message) {

		_, _ = B.Send(m.Sender, fmt.Sprintf("export"))
	})

	B.Handle("/sub", func(m *tb.Message) {

		msg := strings.Split(m.Text, " ")

		if len(msg) == 2 && CheckUrl(msg[1]) {

			url := msg[1]
			registFeed(m.Chat, url)

		} else {
			_, err := B.Send(m.Sender, "请回复RSS URL")
			if err == nil {
				UserState[m.Sender.ID] = fsm.Sub
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

	B.Handle("/set", func(m *tb.Message) {

		sources, _ := model.GetSourcesByUserID(m.Sender.ID)

		var replyButton []tb.ReplyButton
		replyKeys := [][]tb.ReplyButton{}
		for _, source := range sources {
			// 添加按钮
			text := fmt.Sprintf("%s %s", source.Title, source.Link)
			replyButton = []tb.ReplyButton{
				tb.ReplyButton{Text: text},
			}
			replyKeys = append(replyKeys, replyButton)
		}
		_, err := B.Send(m.Sender, "请选择你要设置的源", &tb.ReplyMarkup{
			ReplyKeyboard:   replyKeys,
			OneTimeKeyboard: true,
		})

		if err == nil {
			UserState[m.Sender.ID] = fsm.Set
		}
	})

	B.Handle("/unsub", func(m *tb.Message) {
		msg := strings.Split(m.Text, " ")

		if len(msg) == 2 && CheckUrl(msg[1]) {
			//Unsub by url
			url := msg[1]

			source, _ := model.GetSourceByUrl(url)
			if source == nil {
				_, _ = B.Send(m.Sender, "未订阅该RSS源")
			} else {
				err := model.UnsubByUserIDAndSource(m.Sender.ID, source)
				if err == nil {
					_, _ = B.Send(m.Sender, "退订成功！")
					log.Printf("%d unsubscribe [%d]%s %s", m.Sender.ID, source.ID, source.Title, source.Link)
				} else {
					_, err = B.Send(m.Sender, err.Error())
				}
			}
		} else {
			//Unsub by button
			sources, _ := model.GetSourcesByUserID(m.Sender.ID)
			var replyButton []tb.ReplyButton
			replyKeys := [][]tb.ReplyButton{}
			for _, source := range sources {
				// 添加按钮
				text := fmt.Sprintf("%s %s", source.Title, source.Link)
				replyButton = []tb.ReplyButton{
					tb.ReplyButton{Text: text},
				}

				replyKeys = append(replyKeys, replyButton)
			}
			_, err := B.Send(m.Sender, "请选择你要退订的源", &tb.ReplyMarkup{
				ReplyKeyboard: replyKeys,
			})

			if err == nil {
				UserState[m.Sender.ID] = fsm.UnSub
			}
		}

	})

	B.Handle("/ping", func(m *tb.Message) {

		_, _ = B.Send(m.Sender, "pong")
	})

	B.Handle(tb.OnText, func(m *tb.Message) {
		switch UserState[m.Sender.ID] {
		case fsm.UnSub:
			{
				str := strings.Split(m.Text, " ")
				if len(str) != 2 && !CheckUrl(str[1]) {
					_, _ = B.Send(m.Sender, "请选择正确的指令！")
				} else {
					err := model.UnsubByUserIDAndSourceURL(m.Sender.ID, str[1])
					if err != nil {
						_, _ = B.Send(m.Sender, "请选择正确的指令！")

					} else {
						_, _ = B.Send(
							m.Sender,
							fmt.Sprintf("[%s](%s) 退订成功", str[0], str[1]),
							&tb.SendOptions{
								ParseMode: tb.ModeMarkdown,
							}, &tb.ReplyMarkup{
								ReplyKeyboardRemove: true,
							},
						)
						UserState[m.Sender.ID] = fsm.None
					}
				}
			}

		case fsm.Sub:
			{
				url := strings.Split(m.Text, " ")
				if !CheckUrl(url[0]) {
					_, _ = B.Send(m.Sender, "请回复正确的URL")
					return
				}
				registFeed(m.Chat, url[0])
				UserState[m.Sender.ID] = fsm.None
			}

		case fsm.Set:
			{
				str := strings.Split(m.Text, " ")
				if len(str) != 2 && !CheckUrl(str[1]) {
					_, _ = B.Send(m.Sender, "请选择正确的指令！")
				} else {
					source, err := model.GetSourceByUrl(str[1])

					if err != nil {
						_, _ = B.Send(m.Sender, "请选择正确的指令！")
						return
					}
					sub, err := model.GetSubscribeByUserIDAndSourceID(m.Sender.ID, source.ID)
					if err != nil {
						_, _ = B.Send(m.Sender, "请选择正确的指令！")
						return
					}
					t := template.New("setting template")
					t.Parse(`
订阅<b>设置</b>
[id] {{ .sub.ID}}
[标题] {{ .source.Title }}
[Link] {{.source.Link}}
[通知] {{if eq .sub.EnableNotification 0}}关闭{{else if eq .sub.EnableNotification 1}}开启{{end}}
[Telegraph] {{if eq .sub.EnableTelegraph 0}}关闭{{else if eq .sub.EnableTelegraph 1}}开启{{end}}
`)

					message := new(bytes.Buffer)
					if sub.EnableNotification == 1 {

						toggleNoticeKey.Text = "关闭通知"
					} else {
						toggleNoticeKey.Text = "开启通知"

					}
					if sub.EnableTelegraph == 1 {
						toggleTelegraphKey.Text = "关闭 Telegraph 转码"
					} else {
						toggleTelegraphKey.Text = "开启 Telegraph 转码"
					}
					feedSettingKeys = append(feedSettingKeys, []tb.InlineButton{toggleNoticeKey, toggleTelegraphKey})
					_ = t.Execute(message, map[string]interface{}{"source": source, "sub": sub})
					_, _ = B.Send(
						m.Sender,
						message.String(),
						&tb.SendOptions{
							ParseMode: tb.ModeHTML,
						}, &tb.ReplyMarkup{
							InlineKeyboard: feedSettingKeys,
						},
					)
					UserState[m.Sender.ID] = fsm.None

				}
			}
		}
	})
}
