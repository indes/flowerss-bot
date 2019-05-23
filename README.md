# RSSFlow

[![Build Status](https://travis-ci.org/indes/rssflow.svg?branch=master)](https://travis-ci.org/indes/rssflow)
[![Go Report Card](https://goreportcard.com/badge/github.com/indes/rssflow)](https://goreportcard.com/report/github.com/indes/rssflow)
[![MIT license](https://img.shields.io/github/license/indes/rssflow.svg)](https://github.com/indes/rssflow/blob/master/LICENSE)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Findes%2Frssflow.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Findes%2Frssflow?ref=badge_shield)

DEMO: https://t.me/rssflowbot

<img src="https://raw.githubusercontent.com/indes/rssflow/master/images/rssflow_demo.gif" width = "300"/>


## Features  

- 支持 Telegram 应用内 instant view
- 默认10分钟抓取一次

## 安装

**由于 GoReleaser 不支持 Cgo，如果要使用 SQLite 做为数据库，请下载源码自行编译。**  

### 源码安装

```shell
git clone https://github.com/indes/rssflow && cd rssflow
make build
./rssflow
```

### 下载二进制

**不支持 SQLite**  

从[Releases](https://github.com/indes/rssflow/releases) 页面下载对应的版本。

## 配置

根据以下模板，新建 `config.yml` 文件。

```yml
bot_token: XXX
telegraph_toke: xxxx
socks5: 127.0.0.1:1080
update_interval: 10
mysql:
  host: 123.123.132.132
  port: 3306
  user: user
  password: pwd
  database: rssflow
```

配置说明：

| 配置项 | 含义 | 必填 |
| ------ | ------ | ------ |
| bot_token | Telegram Bot Token | 必填 |
| telegraph_token | Telegraph Token, 用于转存原文到 Telegraph | 可忽略（不转存原文到Telegraph ） |
| update_interval | RSS 源扫描间隔（分钟） | 可忽略（默认10） |
| socks5 | 用于无法正常 Telegram API 的环境 | 可忽略（能正常连接上 Telegram API 服务器） |
| mysql | 数据库配置 | 可忽略（使用 SQLite ） |

## 使用

命令：
```shell
/sub [url] 订阅源
/unsub [url] 取消订阅
/list 查看当前订阅源
/ping :)
```
建议配合 [RSSHub](https://github.com/DIYgod/RSSHub) 使用。

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Findes%2Frssflow.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Findes%2Frssflow?ref=badge_large)
