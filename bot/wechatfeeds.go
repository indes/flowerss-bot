package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/indes/flowerss-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Pagination struct {
	Offset  int
	Limit   int
	Keyword string
	Source  string
}

func (p *Pagination) getPrePagination() Pagination {
	var offset int
	if p.Offset-p.Limit > 0 {
		offset = p.Offset - p.Limit
	} else {
		offset = 0
	}
	return Pagination{
		Offset:  offset,
		Limit:   p.Limit,
		Keyword: p.Keyword,
		Source:  p.Source,
	}
}

func (p *Pagination) getNextPagination() Pagination {
	return Pagination{
		Offset:  p.Offset + p.Limit,
		Limit:   p.Limit,
		Keyword: p.Keyword,
		Source:  p.Source,
	}
}

func searchWechat(keyword string, offset, limit int) (string, *tb.ReplyMarkup) {
	if keyword == "" {
		return "请输入搜索关键字", &tb.ReplyMarkup{}
	}
	pagination := Pagination{
		Offset:  offset,
		Limit:   limit,
		Keyword: keyword,
		Source:  "wechat",
	}
	feeds, total := model.SearchWechatAccounts(keyword, pagination.Offset, pagination.Limit)
	msg := fmt.Sprintf("找到如下 %d 个账号，点击订阅", total)
	fmt.Println(msg)

	// var replyMarkup tb.ReplyMarkup
	var inlineButtons [][]tb.InlineButton
	for _, feed := range feeds {
		// "/sub " + fmt.Sprintf(model.WechatSubUrl, feed.Bizid),
		button := tb.InlineButton{
			Unique: "sub_wechat_feed_item_btn",
			Text:   fmt.Sprintf("「%s」", feed.Name),
			Data:   feed.Bizid,
		}
		inlineButtons = append(inlineButtons, []tb.InlineButton{button})
	}

	pre_btn := tb.InlineButton{}
	next_btn := tb.InlineButton{}
	if pagination.Offset >= pagination.Limit {
		pre_data := pagination.getPrePagination()
		pre_btn = tb.InlineButton{Unique: "pagination_btn", Text: "上一页", Data: fmt.Sprintf("%s:%s:%d:%d", pre_data.Source, pre_data.Keyword, pre_data.Offset, pre_data.Limit)}
	}

	if total > int64(pagination.Limit) {
		next_data := pagination.getNextPagination()
		next_btn = tb.InlineButton{Unique: "pagination_btn", Text: "下一页", Data: fmt.Sprintf("%s:%s:%d:%d", next_data.Source, next_data.Keyword, next_data.Offset, next_data.Limit)}
	}

	inlineButtons = append(inlineButtons, []tb.InlineButton{
		pre_btn,
		next_btn,
	})

	replyMarkup := tb.ReplyMarkup{
		InlineKeyboard: inlineButtons,
	}

	return msg, &replyMarkup
}

func searchCmdCtr(m *tb.Message) {
	keyword := m.Payload
	msg, replyMarkup := searchWechat(keyword, 0, 5)

	_, err := B.Send(m.Chat, msg, replyMarkup)
	if err != nil {
		log.Fatalln("Send Error. ", err)
	}
}

func subWechatFeedItemBtnCtr(c *tb.Callback) {
	bizid := c.Data
	if bizid != "" {
		url := fmt.Sprintf(model.WechatSubUrl, bizid)
		registFeed(c.Message.Chat, url)
		return
	}
	_, _ = B.Edit(c.Message, "订阅错误！")
}

func paginationBtnCtr(c *tb.Callback) {
	data := c.Data
	p_data := strings.Split(data, ":")
	if len(p_data) == 4 {
		offset, err := strconv.Atoi(p_data[2])
		if err != nil {
			log.Fatalln("Get pagination data error", err)
			_, _ = B.Send(c.Message.Chat, "翻页错误！")
		}
		limit, err := strconv.Atoi(p_data[3])
		if err != nil {
			log.Fatalln("Get pagination data error", err)
			_, _ = B.Send(c.Message.Chat, "翻页错误！")
		}
		p := Pagination{
			Offset:  offset,
			Limit:   limit,
			Keyword: p_data[1],
			Source:  p_data[0],
		}
		if p.Source == "wechat" {
			msg, replyMarkup := searchWechat(p.Keyword, p.Offset, p.Limit)
			_, err := B.Edit(c.Message, msg, replyMarkup)
			if err != nil {
				log.Fatalln("Send Error. ", err)
			}
			return
		}
	}

	_, _ = B.Send(c.Message.Chat, "翻页错误！")
}
