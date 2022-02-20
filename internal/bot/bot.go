package bot

import (
	"time"

	"github.com/indes/flowerss-bot/internal/bot/fsm"
	"github.com/indes/flowerss-bot/internal/bot/handler"
	"github.com/indes/flowerss-bot/internal/bot/middleware"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/util"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
)

var (
	// UserState 用户状态，用于标示当前用户操作所在状态
	UserState map[int64]fsm.UserStatus = make(map[int64]fsm.UserStatus)

	// B bot
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
		URL:     config.TelegramEndpoint,
		Token:   config.BotToken,
		Poller:  spamProtected,
		Client:  util.HttpClient,
		Verbose: true,
	})
	B.Use(middleware.PreLoadMentionChat(), middleware.IsChatAdmin())
	if err != nil {
		zap.S().Fatal(err)
		return
	}
}

//Start bot
func Start() {
	if config.RunMode == config.TestMode {
		return
	}

	zap.S().Infof("bot start %s", config.AppVersionInfo())
	setCommands()
	B.Start()
}

func setCommands() {
	commandHandlers := []handler.CommandHandler{
		&handler.Start{},
		&handler.Version{},
		handler.NewPing(B),
		handler.NewRemoveSubscription(B),
		&handler.Help{},
		&handler.Export{},
		&handler.Import{},
		&handler.AddSubscription{},
		&handler.ListSubscription{},
		&handler.RemoveAllSubscription{},
		handler.NewOnDocument(B),
		handler.NewPauseAll(),
		handler.NewActiveAll(),
		handler.NewSetFeedTag(),
		handler.NewSetUpdateInterval(),
		handler.NewSet(B),
	}

	for _, h := range commandHandlers {
		B.Handle(h.Command(), h.Handle, h.Middlewares()...)
	}

	ButtonHandlers := []handler.ButtonHandler{
		&handler.RemoveAllSubscriptionButton{},
		&handler.CancelRemoveAllSubscriptionButton{},
	}

	for _, h := range ButtonHandlers {
		B.Handle(h, h.Handle, h.Middlewares()...)
	}
	// 设置bot命令提示信息
	//commands := []tb.Command{
	//	{Text: "start", Description: "开始使用"},
	//	{Text: "sub", Description: "订阅rss源"},
	//	{Text: "list", Description: "当前订阅的rss源"},
	//	{Text: "unsub", Description: "退订rss源"},
	//	{Text: "unsuball", Description: "退订所有rss源"},
	//
	//	{Text: "set", Description: "设置rss订阅"},
	//	{Text: "setfeedtag", Description: "设置rss订阅标签"},
	//	{Text: "setinterval", Description: "设置rss订阅抓取间隔"},
	//
	//	{Text: "export", Description: "导出订阅为opml文件"},
	//	{Text: "import", Description: "从opml文件导入订阅"},
	//
	//	{Text: "check", Description: "检查我的rss订阅状态"},
	//	{Text: "pauseall", Description: "停止抓取订阅更新"},
	//	{Text: "activeall", Description: "开启抓取订阅更新"},
	//
	//	{Text: "help", Description: "使用帮助"},
	//	{Text: "version", Description: "bot版本"},
	//}

	var commands []tb.Command
	for _, h := range commandHandlers {
		if h.Description() == "" {
			continue
		}
		commands = append(commands, tb.Command{Text: h.Command(), Description: h.Description()})
	}
	zap.S().Debugf("set bot command %+v", commands)
	if err := B.SetCommands(commands); err != nil {
		zap.S().Errorw("set bot commands failed", "error", err.Error())
	}
}

func setHandle() {
	//B.Handle(&tb.InlineButton{Unique: "set_feed_item_btn"}, setFeedItemBtnCtr)
	//B.Handle(&tb.InlineButton{Unique: "set_toggle_notice_btn"}, setToggleNoticeBtnCtr)
	//B.Handle(&tb.InlineButton{Unique: "set_toggle_telegraph_btn"}, setToggleTelegraphBtnCtr)
	//B.Handle(&tb.InlineButton{Unique: "set_toggle_update_btn"}, setToggleUpdateBtnCtr)
	//B.Handle(&tb.InlineButton{Unique: "set_set_sub_tag_btn"}, setSubTagBtnCtr)
	//B.Handle(&tb.InlineButton{Unique: "unsub_all_cancel_btn"}, unsubAllCancelBtnCtr)
	//B.Handle(&tb.InlineButton{Unique: "unsub_feed_item_btn"}, unsubFeedItemBtnCtr)

	//B.Handle(tb.OnText, textCtr)
}
