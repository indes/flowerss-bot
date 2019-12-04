package bot

import (
	"github.com/indes/flowerss-bot/bot/fsm"
	"github.com/indes/flowerss-bot/config"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"time"
)

var (
	botToken                             = config.BotToken
	socks5Proxy                          = config.Socks5
	UserState   map[int64]fsm.UserStatus = make(map[int64]fsm.UserStatus)

	B              *tb.Bot
)

func init() {
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	spamProtected := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {

		if !CheckAdmin(upd) {
			return false
		}

		return true
	})
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
			Poller: spamProtected,
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
			Poller: spamProtected,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
	}

}

//Start bot
func Start() {
	makeHandle()
	B.Start()
}

func makeHandle() {
	B.Handle(&tb.InlineButton{Unique: "set_toggle_notice_btn"}, setToggleNoticeBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_toggle_telegraph_btn"},setToggleTelegraphBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_toggle_update_btn"}, SetToggleUpdateBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "unsub_all_confirm_btn"}, unsubAllConfirmBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "unsub_all_cancel_btn"}, unsubAllCancelBtnCtr)

	B.Handle("/start", startCmdCtr)

	B.Handle("/export", exportCmdCtr)

	B.Handle("/sub", subCmdCtr)

	B.Handle("/list", listCmdCtr)

	B.Handle("/set", setCmdCtr)

	B.Handle("/unsub", unsubCmdCtr)

	B.Handle("/unsuball", unsubAllCmdCtr)

	B.Handle("/ping", pingCmdCtr)

	B.Handle("/help", helpCmdCtr)

	B.Handle("/import", importCmdCtr)

	B.Handle(tb.OnText, textCtr)

	B.Handle(tb.OnDocument, docCtr)
}
