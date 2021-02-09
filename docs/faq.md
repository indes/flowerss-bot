

### 日志中大量类似于 `Create telegraph page error: FLOOD_WAIT_7` 的提示。  

原因是创建 Telegraph 页面请求过快触发了接口限制，可尝试在配置文件中添加多个 Telegraph token。 


### 如何申请 Telegraph Token？ 

如果要使用应用内即时预览，必须在配置文件中填写 `telegraph_token` 配置项，Telegraph Token 申请命令如下：  
```bash
curl https://api.telegra.ph/createAccount?short_name=flowerss&author_name=flowerss&author_url=https://github.com/indes/flowerss-bot
```

返回的 JSON 中 access_token 字段值即为 Telegraph Token。


### 如何获取我的telegram id？
可以参考这个网页获取：https://botostore.com/c/getmyid_bot/