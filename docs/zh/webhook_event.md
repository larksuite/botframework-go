# webhook-event  
对于开放平台订阅事件  
1. 将 feishu 或 lark 开放平台应用事件订阅配置页面上的请求网址URL与应用后端服务做关联，即该 URL 的请求最终会请求到应用后端服务。
2. 应用后端服务调用 SDK 的 event.EventCallback 函数进行协议解析，将不同的事件路由到对应的事件处理函数中，处理相应业务逻辑。    

本框架支持自动生成基于 gin 框架的应用后端代码，其中已经实现了事件订阅的相关操作，可供开发者使用或参考。
  
## 绑定回调处理函数  
将应用事件订阅配置页面上的请求网址URL指向应用后端服务之后
如果使用 gin 框架，在 main 函数中添加如下代码，其中 EventCallback 为本框架生成的订阅事件回调处理函数。  
```go
r := gin.Default()

r.POST("/webhook/event", EventCallback) //open platform event callback
```
  
## 回调与解析  
如果使用了自动生成的代码，则回调与解析函数也已经自动生成了，也就是生成代码 `callback.go` 文件中的 EventCallback 函数。    
只需要将函数中的 appID 替换为正确的应用 ID 即可。  

下面给出 EventCallback 示例代码。    
```go
// EventCallback open platform event
func EventCallback(c *gin.Context) {
    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil || len(body) == 0 {
        c.JSON(500, gin.H{"codemsg": common.ErrEventParams.String()})
        return
    }
    appID := "your appid"
    challenge, err := event.EventCallback(c, string(body), appID)
    if err != nil {
        c.JSON(500, gin.H{"codemsg": fmt.Sprintf("%v", err)})
    } else if "" != challenge {
        c.JSON(200, gin.H{"challenge": challenge})
    } else {
        c.JSON(200, gin.H{"codemsg": common.Success.String()})
    }
}
```
  
手动编写代码:    
主要是通过调用 SDK/event/event.go 中的 EventCallback 函数进行对 http 请求包中的 body 进行解析。  
  
## 事件注册和处理  
本框架是基于注册和回调函数来实现事件的分发和处理的。  
开发者需要先注册回调函数，之后 event.EventCallback 函数将解析请求协议并调用对应事件的回调处理函数，来执行对应的业务逻辑。  

框架自动生成代码:      
关于事件注册，框架会基于配置文件自动生成注册函数和回调处理函数。    
在生成代码的 handler_event/regist.go 文件中有 RegistHandler 函数，函数里注册了一系列相应的事件处理函数。    
简单的示例代码如下:    
```go
func RegistHandler(appID string) {
    // regist open platform event handler
    event.EventRegister(appID, protocol.EventTypeMessage, EventMessage)
    event.EventRegister(appID, protocol.EventTypeAppTicket, EventAppTicket)
    event.EventRegister(appID, protocol.EventTypeAppOpen, EventAppOpen)
    event.EventRegister(appID, protocol.EventTypeAddBot, EventAddBot)
    event.EventRegister(appID, protocol.EventTypeP2PChatCreate, EventP2PChatCreate)

    // regist bot recv message handler
    event.BotRecvMsgRegister(appID, "default", BotRecvMsgDefault)
    event.BotRecvMsgRegister(appID, "help", BotRecvMsgHelp)

}
```
event.EventRegister 函数用于注册不同种类的订阅事件。  
event.BotRecvMsgRegister 函数用于注册机器人接收消息时不同消息首单词的回调函数，它是 protocol.EventTypeMessage 类型事件的细分，方便基于接收消息首单词执行不同的业务操作。  
event.BotRecvMsgRegister 函数中第二个参数为关键字，当消息首单词为该关键词时，将指向对应的回调函数，其中当未匹配到任何关键字时，将调用 default 关键字对应的回调函数。  

**用户和机器人会话首次被创建事件回调函数，示例代码如下**           
```go
//EventP2PChatCreate在会话首次被创建时向用户发送"hello,P2PChatCreate"消息
func EventP2PChatCreate(ctx context.Context, eventBody []byte) error {
    request := &protocol.P2PChatCreateEvent{}
    err := json.Unmarshal(eventBody, request)
    if err != nil {
        return err
    }
    user:=&protocol.UserInfo{
        ID:   request.ChatID,
        Type: protocol.UserTypeChatID,
    }
    message.SendTextMessage(ctx,"your tenanKey",request.AppID,user,"","hello,P2PChatCreate")
    return nil
}
``` 
  
**ISVApp 的 app_ticket 事件回调函数，示例代码如下**      
```go
//EventAppTicket函数自动更新AppTicket
func EventAppTicket(ctx context.Context, eventBody []byte) error {
    return auth.RefreshAppTicket(ctx, eventBody)
}
```

**BotRecvMsgDefault 函数示例代码如下**  
```go
//（该函数会向用户回复"Text that is empty or is not matched"消息）
func BotRecvMsgDefault(ctx context.Context, msg *protocol.BotRecvMsg) error {
   user:=&protocol.UserInfo{
      ID:   msg.OpenChatID,
      Type: protocol.UserTypeChatID,
   }
   input := msg.TextParam
   message.SendTextMessage(ctx,"your tenantKey",msg.AppID,user,"","Text that is empty or is not matched")
   return nil
}
```
  
**BotRecvMsgHelp 函数示例代码如下**   
```go
//（该函数会向用户回复"this is help"消息）
func BotRecvMsgHelp(ctx context.Context, msg *protocol.BotRecvMsg) error {
   user:=&protocol.UserInfo{
      ID:   msg.OpenChatID,
      Type: protocol.UserTypeChatID,
   }
   message.SendTextMessage(ctx,"your tenantKey",msg.AppID,user,"","this is help")
   return nil
}
```
  
手动编写代码:    
利用 event.EventRegister 和 event.BotRecvMsgRegister 函数注册需要处理的订阅事件即可。    
你可以参考 gin 框架自动生成代码 中的 RegistHandler 函数去实现它。  
