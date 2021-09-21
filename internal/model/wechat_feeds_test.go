package model

import "testing"
import "github.com/stretchr/testify/assert"

func Test_ProcessWechatURL(t *testing.T) {
	assert := assert.New(t)

	// Test normal url
	urlForSub := "https://github.com/hellodword/wechat-feeds/raw/feeds/Mzg2ODAyNTgyMQ==.xml"
	res := ProcessWechatURL(urlForSub)
	assert.Equal(res, urlForSub)

	// Test wechat url with __biz: Success
	url := "https://mp.weixin.qq.com/s?__biz=Mzg2ODAyNTgyMQ==&mid=2247484843&idx=1&sn=be4445bf3a7520b86d15876d46b7b8ab&chksm=ceb3d119f9c4580f6b360299affd4e2d312a619d32edc2a5e1c031f3a4c600b4b52132d5eb56&scene=132#wechat_redirect"
	res = ProcessWechatURL(url)
	assert.Equal(res, urlForSub)

	// Test wechat url without __biz: Like normal url
	url = "https://mp.weixin.qq.com/s/nGTr26rhHIUJqq6bV8Xwsw"
	res = ProcessWechatURL(url)
	assert.Equal(res, url)

	// Test invalid url
	url = "abcd"
	res = ProcessWechatURL(url)
	assert.Equal(res, url)
}
