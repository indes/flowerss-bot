# RSSFlow

[![Go Report Card](https://goreportcard.com/badge/github.com/indes/rssflow)](https://goreportcard.com/report/github.com/indes/rssflow)
[![MIT license](https://img.shields.io/github/license/indes/rssflow.svg)](https://github.com/indes/rssflow/blob/master/LICENSE)

DEMO: https://t.me/rssflowbot

## Features  

- 支持 Telegram 应用内 instant view
- 默认10分钟抓取一次

## 安装

**由于 GoReleaser 不支持 Cgo，如果要使用 SQLite 做为数据库，请下载源码自行编译。**  

### 源码安装

```shell
git clone https://github.com/indes/rssflow && cd rssflow
go model download
go build .
./rssflow
```

### 下载二进制

**不支持 SQLite**  

从[Releases](https://github.com/indes/rssflow/releases) 页面下载对应的版本。

## 配置

根据以下模板，新建 `config.yml` 文件。

```yml
token: XXX
socks5: XXX
mysql:
  host: XXX
  port: XXX
  user: XXX
  password: XXX
  database: XXX
```

配置说明：

| 配置项 | 含义 | 必填 |
| ------ | ------ | ------ |
| token | Telegram Bot Token | 必填 |
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