package bot

import (
	"time"

	"github.com/indes/flowerss-bot/internal/bot/fsm"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/util"
	
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	// UserState 用户状态，用于标示当前用户操作所在状态
	UserState map[int64]fsm.UserStatus = make(map[int64]fsm.UserStatus)

	// B telebot
	B *tb.Bot
)

func init() {
	if config.RunMode == config.TestMode {
		return
	}
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	spamProtected := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		if !isUserAllowed(upd) {
			// 检查用户是否可以使用bot
			return false
		}

		if !CheckAdmin(upd) {
			return false
		}
		return true
	})
	zap.S().Infow("init telegram bot",
		"token", config.BotToken,
		"endpoint", config.TelegramEndpoint,
	)

	// create bot
	var err error

	B, err = tb.NewBot(tb.Settings{
		URL:    config.TelegramEndpoint,
		Token:  config.BotToken,
		Poller: spamProtected,
		Client: util.HttpClient,
	})

	if err != nil {
		zap.S().Fatal(err)
		return
	}
}

//Start bot
func Start() {
	if config.RunMode != config.TestMode {
		zap.S().Infof("bot start %s", config.AppVersionInfo())
		setCommands()
		setHandle()
		B.Start()
	}
}

func setCommands() {
	// 設置bot命令提示信息
	commands := []tb.Command{
{Text: "start", Description: "開始使用"},
		{Text: "sub", Description: "訂閱rss源"},
		{Text: "list", Description: "當前訂閱的rss源"},
		{Text: "unsub", Description: "退訂rss源"},
		{Text: "unsuball", Description: "退訂所有rss源"},

		{Text: "set", Description: "設置rss訂閱"},
		{Text: "setfeedtag", Description: "設置rss訂閱標籤"},
		{Text: "setinterval", Description: "設置rss訂閱抓取間隔"},

		{Text: "export", Description: "導出訂閱為opml文件"},
		{Text: "import", Description: "從opml文件導入訂閱"},

		{Text: "check", Description: "檢查我的rss訂閱狀態"},
		{Text: "pauseall", Description: "停止抓取訂閱更新"},
		{Text: "activeall", Description: "開啟抓取訂閱更新"},

		{Text: "help", Description: "使用幫助"},
		{Text: "version", Description: "bot版本"},
	}

	zap.S().Debugf("set bot command %+v", commands)

	if err := B.SetCommands(commands); err != nil {
		zap.S().Errorw("set bot commands failed", "error", err.Error())
	}
}

func setHandle() {
	B.Handle(&tb.InlineButton{Unique: "set_feed_item_btn"}, setFeedItemBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_toggle_notice_btn"}, setToggleNoticeBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_toggle_telegraph_btn"}, setToggleTelegraphBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_toggle_update_btn"}, setToggleUpdateBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_set_sub_tag_btn"}, setSubTagBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "unsub_all_confirm_btn"}, unsubAllConfirmBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "unsub_all_cancel_btn"}, unsubAllCancelBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "unsub_feed_item_btn"}, unsubFeedItemBtnCtr)

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

	B.Handle("/setfeedtag", setFeedTagCmdCtr)

	B.Handle("/setinterval", setIntervalCmdCtr)

	B.Handle("/check", checkCmdCtr)

	B.Handle("/activeall", activeAllCmdCtr)

	B.Handle("/pauseall", pauseAllCmdCtr)

	B.Handle("/version", versionCmdCtr)

	B.Handle(tb.OnText, textCtr)

	B.Handle(tb.OnDocument, docCtr)
}
