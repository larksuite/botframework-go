# webhook-event    

## bind  
generated code:    
Code for binding has been generated in main.go.    
```go
r := gin.Default()

r.POST("/webhook/event", EventCallback) //open platform event callback
```
   
## call back  
generated code:    
Code for event callback has been generated in `callback.go`.    
You can use it just by providing some params.    
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
  
your own code:    
Use EventCallback function in SDK/event/event.go to analyze http request body.    
  
## register  
generated code:   
You can find RegistHandler function which was in regist.go. There are some response functions.   
 
Register Example code:  
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

**create chat firstly example code**:  
```go
//add EventP2PChatCreate function
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
  
**app_ticket for ISVApp example code**:  
```go
//register AppTicket event
func RegistHandler(appID string) {
   event.EventRegister(appID, protocol.EventTypeAppTicket, EventAppTicket)
}
//add EventAppTicket function
func EventAppTicket(ctx context.Context, eventBody []byte) error {
   return auth.RefreshAppTicket(ctx, eventBody)
}
```
  
**BotRecvMsgDefault example code**:  
```go
//send "text that is empty or is not matched"
func BotRecvMsgDefault(ctx context.Context, msg *protocol.BotRecvMsg) error {
   user:=&protocol.UserInfo{
      ID:   msg.OpenChatID,
      Type: protocol.UserTypeChatID,
   }
   input := msg.TextParam
   message.SendTextMessage(ctx,"your tenantKey",msg.AppID,user,"","Text that is empty or isnot matched")
   return nil
}
```
  
**BotRecvMsgHelp example code**:  
```go
//send "this is help"
func BotRecvMsgHelp(ctx context.Context, msg *protocol.BotRecvMsg) error {
   user:=&protocol.UserInfo{
      ID:   msg.OpenChatID,
      Type: protocol.UserTypeChatID,
   }
   message.SendTextMessage(ctx,"your tenantKey",msg.AppID,user,"","this is help")
   return nil
}
```  
