package model

import (
	"fmt"
	"net/url"
)
const wechatHost = "mp.weixin.qq.com"
const wechatSubUrl = "https://github.com/hellodword/wechat-feeds/raw/feeds/%s.xml"

// ProcessWechatURL return wechat-feed sub URL. if it's valid wehchat url, else return origin str
func ProcessWechatURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err == nil {
		if (u.Host == "mp.weixin.qq.com") {
			q := u.Query()
			bizs, ok := q["__biz"]
			if ok {
				biz := bizs[0]
				newURL := fmt.Sprintf(wechatSubUrl, biz)
				return newURL
			}
		}
	}
	return urlStr
}