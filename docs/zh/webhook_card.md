# webhook-card  
对于开放平台卡片交互回调  
1. 在 feishu 或 lark 开放平台应用信息页面，将应用功能--机器人选项卡配置页面上的消息卡片请求网址 URL 与应用后端服务做关联，即该 URL 的请求最终会请求到应用后端服务。  
2. 应用后端服务调用 SDK 的 event.CardCallBack 函数进行协议解析，将不同的交互操作路由到对应的回调处理函数中，处理相应业务逻辑。    

本框架支持自动生成基于 gin 框架的应用后端代码，其中已经实现了卡片交互回调的注册和调用操作，可供开发者使用或参考。  

## 绑定  
将消息卡片请求网址 URL 指向应用后端服务之后  
如果使用 gin 框架，在 main 函数中添加如下代码，其中 CardCallback 为本框架生成的订阅事件回调处理函数。  
```go
r := gin.Default()

r.POST("/webhook/card", CardCallback)   //card action callback
```

## 回调与解析  
如果使用了自动生成的代码，则回调与解析函数也已经自动生成了，也就是生成代码 callback.go 文件中的 CardCallback 函数。    
只需要将函数中的 appID 替换为正确的应用 ID 即可。  

下面给出 CardCallback 示例代码。   
```go
// CardCallback card action callback
func CardCallback(c *gin.Context) {
   body, err := ioutil.ReadAll(c.Request.Body)
   if err != nil || len(body) == 0 {
      c.JSON(500, gin.H{"codemsg": common.ErrCardParams.String()})
      return
   }
   // for verify signature
   header := map[string]string{
      "X-Lark-Request-Timestamp": c.Request.Header.Get("X-Lark-Request-Timestamp"),
      "X-Lark-Request-Nonce":     c.Request.Header.Get("X-Lark-Request-Nonce"),
      "X-Lark-Signature":         c.Request.Header.Get("X-Lark-Signature"),
   }
   appID := "your appid"
   card, challenge, err := event.CardCallBack(c, appID, header, body)
   if err != nil {
      c.JSON(500, gin.H{"codemsg": fmt.Sprintf("%v", err)})
   } else if "" != challenge {
      c.JSON(200, gin.H{"challenge": challenge})
   } else {
      data, err := json.Marshal(card)
      if err != nil {
         c.JSON(500, gin.H{"codemsg": fmt.Sprintf("%v", err)})
      } else {
         c.String(200, string(data))
      }
   }
}
```
  
## 事件注册和处理  
用户点击卡片上的交互模块之后，应用后端服务将接收到对应的请求，框架将依据注册信息，路由到对应的回调函数执行业务逻辑。    

以“机器人接收以 card 为首单词的消息后，后端服务向用户发送卡片消息，用户收到卡片消息后点击卡片上的按钮，后端服务接收用户点击操作的信息”这个业务逻辑为例，说明卡片交互事件的处理流程。  

1. 注册 event 事件。  
```go
func RegistHandler(appID string) {
   //机器人接收消息回调注册
   event.EventRegister(appID, protocol.EventTypeMessage, EventMessage)

   //以 card 为首单词的消息处理回调函数注册
   event.BotRecvMsgRegister(appID, "card", BotRecvMsgCard)

   //卡片交互回调注册：将method   clickbutton   和函数ActionClickButton绑定
   event.CardRegister(appID, "clickbutton", ActionClickButton)
}
```
  
2. 补全 EventMessage 函数和 BotRecvMsgCard 函数。    
（如果使用gin框架自动生成的代码则 EventMessage 函数已经生成了）  
```go
func EventMessage(ctx context.Context, eventBody []byte) error {
   return event.BotRecvMsgHandler(ctx, eventBody)
}

func BotRecvMsgCard(ctx context.Context, msg *protocol.BotRecvMsg) error {
   user := &protocol.UserInfo{
      ID:   msg.OpenChatID,
      Type: protocol.UserTypeChatID,
   }

   //build card
   builder := &message.CardBuilder{}
   //add config
   config := protocol.ConfigForm{
      MinVersion:     protocol.VersionForm{},
      WideScreenMode: true,
   }
   builder.SetConfig(config)

   //add header
   content := "this is a card"
   line := 1
   title := protocol.TextForm{
      Tag:     protocol.PLAIN_TEXT_E,
      Content: &content,
      Lines:   &line,
   }
   builder.AddHeader(title, "")

   //add button
   button1 := make(map[string]string, 0)
   button1["key"] = "buttonValue"
   builder.AddActionBlock([]protocol.ActionElement{
      message.NewButton(message.NewMDText("callback button", nil, nil, nil), nil, nil, button1, protocol.DANGER, nil, "clickbutton"),
   })

   card, err := builder.BuildForm()
   if err != nil {
      return fmt.Errorf("card build failed error[%v]", err)
   }

   //add params to use message.SendCardMessage
   //最后一个参数为卡片的类型，独享或共享。当为false时，卡片为独享卡片，只会更新指定用户看到的卡片(发出更新卡片请求的用户的卡片一定会被更新)。未指定用户的卡片并不会被更新。而当updateMulti为true时，卡片为共享卡片，所有用户的卡片都会被同步更新。
   _, err := message.SendCardMessage(ctx, tenantKey, appID, user, "", *card, false)
   if err != nil {
      return fmt.Errorf("send message failed error[%v]", err)
   }

   return nil
}
```
  
上面的 BotRecvMsgAsyncButton 函数会在用户发送"card"后，回复一张含有一个名为"callback button"的按钮的卡片。    
点击消息卡片中的按钮，便会发送相应的消息卡片请求。    
    
3. 实现响应函数    
```go
func ActionClickButton(ctx context.Context, callback *protocol.CardCallbackForm) (*protocol.CardForm, error) {
   method, _ := callback.Action.Value["method"]
   sessionID, _ := callback.Action.Value["sid"]
   common.Logger(ctx).Infof("cardActionCallBack: method[%s]sessionID[%s]", method, sessionID)

   //build card
   builder := &message.CardBuilder{}

   //add openids  (如果你想更新其他用户的卡片，可以向OpenIDs参数中加入对应用户的openid，如果不需要，可以不写)
   openids := make([]string,1)
   openids[0] = "another user's openid"
   builder.OpenIDs = openids

   //add config
   config := protocol.ConfigForm{
      MinVersion:     protocol.VersionForm{},
      WideScreenMode: true,
   }
   builder.SetConfig(config)

   //add header
   content := "this is a card"
   line := 1
   title := protocol.TextForm{
      Tag:     protocol.PLAIN_TEXT_E,
      Content: &content,
      Lines:   &line,
   }
   builder.AddHeader(title, "")

   //add button
   button1 := make(map[string]string, 0)
   button1["key"] = "buttonValue"
   builder.AddActionBlock([]protocol.ActionElement{
      message.NewButton(message.NewMDText("clicked button", nil, nil, nil), nil, nil, button1, protocol.UNKNOWN, nil, "clickbutton"),
   })

   card, err := builder.BuildForm()
   if err != nil {
      return &protocol.CardForm{}, fmt.Errorf("card update failed error[%v]", err)
   }

   return card, nil
}
```  
该响应函数实现: 用户点击按钮后，将卡片上的按钮更新为灰色不可点击状态。
