# 無料案內所RSS-BOT

[![Build Status](https://github.com/makubex2010/flowerss-bot/workflows/Release/badge.svg)](https://github.com/makubex2010/flowerss-bot/actions?query=workflow%3ARelease)
[![Test Status](https://github.com/makubex2010/flowerss-bot/workflows/Test/badge.svg)](https://github.com/makubex2010/flowerss-bot/actions?query=workflow%3ATest)
[![Build Docker Image](https://github.com/makubex2010//flowerss-bot/workflows/Build%20Docker%20Image/badge.svg)](https://github.com/makubex2010/flowerss-bot/actions?query=workflow%3Adocker)
[![Go Report Card](https://goreportcard.com/badge/github.com/makubex2010/flowerss-bot)](https://goreportcard.com/report/github.com/makubex2010/flowerss-bot)
[![GitHub](https://img.shields.io/github/license/makubex2010/flowerss-bot.svg)](https://github.com/makubex2010/flowerss-bot/blob/master/LICENSE)

[安裝與使用文檔](https://github.com/makubex2010/RSS-BOT)  

<img src="https://github.com/rssflow/img/raw/master/images/rssflow_demo.gif" width = "300"/>

## 功能

- 常見的 RSS Bot 該有的功能
- 支持 Telegram 應用內 instant view
- 支持為 Group 和 Channel 訂閱 RSS 消息
- 豐富的訂閱設置

## 安裝與使用

詳細安裝與使用方法請查閱項目[使用文檔](https://github.com/makubex2010/RSS-BOT)。

使用命令：

```
/sub [url] 訂閱（url 為可選）
/unsub [url] 取消訂閱（url 為可選）
/list 查看當前訂閱
/set 設置訂閱
/check 檢查當前訂閱
/setfeedtag [sub id] [tag1] [tag2] 設置訂閱標籤（最多設置三個Tag，以空格分割）
/setinterval [interval] [sub id] 設置訂閱刷新頻率（可設置多個sub id，以空格分割）
/activeall 開啟所有訂閱
/pauseall 暫停所有訂閱
/import 導入 OPML 文件
/export 導出 OPML 文件
/unsuball 取消所有訂閱
/help 幫助
```
詳細使用方法請查閱項目[使用文檔](https://github.com/makubex2010/RSS-BOT)。
