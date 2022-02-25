package handler

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/model"
)

type Set struct {
	bot *tb.Bot
}

func NewSet(bot *tb.Bot) *Set {
	return &Set{bot: bot}
}

func (s *Set) Command() string {
	return "/set"
}

func (s *Set) Description() string {
	return "设置订阅"
}

func (s *Set) Handle(ctx tb.Context) error {
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	ownerID := ctx.Message().Chat.ID
	if mentionChat != nil {
		ownerID = mentionChat.ID
	}

	sources, err := model.GetSourcesByUserID(ownerID)
	if err != nil {
		return ctx.Reply("获取订阅失败")
	}
	if len(sources) <= 0 {
		return ctx.Reply("当前没有订阅")
	}

	// 配置按钮
	var replyButton []tb.ReplyButton
	replyKeys := [][]tb.ReplyButton{}
	setFeedItemBtns := [][]tb.InlineButton{}
	for _, source := range sources {
		// 添加按钮
		text := fmt.Sprintf("%s %s", source.Title, source.Link)
		replyButton = []tb.ReplyButton{
			tb.ReplyButton{Text: text},
		}
		replyKeys = append(replyKeys, replyButton)

		setFeedItemBtns = append(
			setFeedItemBtns, []tb.InlineButton{
				tb.InlineButton{
					Unique: SetFeedItemButtonUnique,
					Text:   fmt.Sprintf("[%d] %s", source.ID, source.Title),
					Data:   fmt.Sprintf("%d:%d", ownerID, source.ID),
				},
			},
		)
	}

	return ctx.Reply(
		"请选择你要设置的源", &tb.ReplyMarkup{
			InlineKeyboard: setFeedItemBtns,
		},
	)
}

func (s *Set) Middlewares() []tb.MiddlewareFunc {
	return nil
}

const (
	SetFeedItemButtonUnique = "set_feed_item_btn"
	feedSettingTmpl         = `
订阅<b>设置</b>
[id] {{ .sub.ID }}
[标题] {{ .source.Title }}
[Link] {{.source.Link }}
[抓取更新] {{if ge .source.ErrorCount .Count }}暂停{{else if lt .source.ErrorCount .Count }}抓取中{{end}}
[抓取频率] {{ .sub.Interval }}分钟
[通知] {{if eq .sub.EnableNotification 0}}关闭{{else if eq .sub.EnableNotification 1}}开启{{end}}
[Telegraph] {{if eq .sub.EnableTelegraph 0}}关闭{{else if eq .sub.EnableTelegraph 1}}开启{{end}}
[Tag] {{if .sub.Tag}}{{ .sub.Tag }}{{else}}无{{end}}
`
)

type SetFeedItemButton struct {
	bot *tb.Bot
}

func NewSetFeedItemButton(bot *tb.Bot) *SetFeedItemButton {
	return &SetFeedItemButton{bot: bot}
}

func (r *SetFeedItemButton) CallbackUnique() string {
	return "\f" + SetFeedItemButtonUnique
}

func (r *SetFeedItemButton) Description() string {
	return ""
}

func (r *SetFeedItemButton) Handle(ctx tb.Context) error {
	data := strings.Split(ctx.Callback().Data, ":")
	if len(data) < 2 {
		return nil
	}
	subscriberID, _ := strconv.ParseInt(data[0], 10, 64)
	// 如果订阅者与按钮点击者id不一致，需要验证管理员权限
	if subscriberID != ctx.Callback().Sender.ID {
		channelChat, err := r.bot.ChatByUsername(fmt.Sprintf("%d", subscriberID))
		if err != nil {
			return ctx.Edit("获取订阅信息失败")
		}

		if !chat.IsChatAdmin(r.bot, channelChat, ctx.Callback().Sender.ID) {
			return ctx.Edit("获取订阅信息失败")
		}
	}

	sourceID, _ := strconv.Atoi(data[1])
	source, err := model.GetSourceById(uint(sourceID))
	if err != nil {
		return ctx.Edit("找不到该订阅源")
	}

	sub, err := model.GetSubscribeByUserIDAndSourceID(subscriberID, source.ID)
	if err != nil {
		return ctx.Edit("用户未订阅该rss")
	}

	t := template.New("setting template")
	_, _ = t.Parse(feedSettingTmpl)
	text := new(bytes.Buffer)
	_ = t.Execute(text, map[string]interface{}{"source": source, "sub": sub, "Count": config.ErrorThreshold})
	return ctx.Edit(
		text.String(),
		&tb.SendOptions{ParseMode: tb.ModeHTML},
		&tb.ReplyMarkup{InlineKeyboard: genFeedSetBtn(ctx.Callback(), sub, source)},
	)
}

func genFeedSetBtn(
	c *tb.Callback, sub *model.Subscribe, source *model.Source,
) [][]tb.InlineButton {
	setSubTagKey := tb.InlineButton{
		Unique: "set_set_sub_tag_btn",
		Text:   "标签设置",
		Data:   c.Data,
	}

	toggleNoticeKey := tb.InlineButton{
		Unique: "set_toggle_notice_btn",
		Text:   "开启通知",
		Data:   c.Data,
	}
	if sub.EnableNotification == 1 {
		toggleNoticeKey.Text = "关闭通知"
	}

	toggleTelegraphKey := tb.InlineButton{
		Unique: TelegraphSwitchButtonUnique,
		Text:   "开启 Telegraph 转码",
		Data:   c.Data,
	}
	if sub.EnableTelegraph == 1 {
		toggleTelegraphKey.Text = "关闭 Telegraph 转码"
	}

	toggleEnabledKey := tb.InlineButton{
		Unique: SubscriptionSwitchButtonUnique,
		Text:   "暂停更新",
		Data:   c.Data,
	}

	if source.ErrorCount >= config.ErrorThreshold {
		toggleEnabledKey.Text = "重启更新"
	}

	feedSettingKeys := [][]tb.InlineButton{
		[]tb.InlineButton{
			toggleEnabledKey,
			toggleNoticeKey,
		},
		[]tb.InlineButton{
			toggleTelegraphKey,
			setSubTagKey,
		},
	}
	return feedSettingKeys
}

func (r *SetFeedItemButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}
