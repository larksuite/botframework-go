# webhook-card  

## bind  
generated code:    
Code for binding has been generated in main.go.  
```go
r := gin.Default()

r.POST("/webhook/card", CardCallback)   //card action callback
```      

## callback  
generated code:    
Code for card callback has been generated in callback.go. You can use it just by providing some params.  
Example code:  
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


## register and callback
Card callback functions were be registered with interaction modules. Here we use code to explain in detail.  
Example code:  
1. registe event callback
```go
func RegistHandler(appID string) {
   event.EventRegister(appID, protocol.EventTypeMessage, EventMessage)

   event.BotRecvMsgRegister(appID, "card", BotRecvMsgCard)

   event.CardRegister(appID, "clickbutton", ActionClickButton)
}
```
  
2. add BotRecvMsgCard function  (send a card)  
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
   EnableForward: false,
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
   //last param means updateMulti. If the param is true, the card type is shared, all users have the same card.
   //If the param is false, the card type is exclusive, usre can only change his own card in general.
   _, err := message.SendCardMessage(ctx, tenantKey, appID, user, "", *card, false)
   if err != nil {
      return fmt.Errorf("send message failed error[%v]", err)
   }

   return nil
}
```
  
3. add ActionClickButton function(update the card)  
```go
func ActionClickButton(ctx context.Context, callback *protocol.CardCallbackForm) (*protocol.CardForm, error) {
   method, _ := callback.Action.Value["method"]
   sessionID, _ := callback.Action.Value["sid"]
   common.Logger(ctx).Infof("cardActionCallBack: method[%s]sessionID[%s]", method, sessionID)

   //build card
   builder := &message.CardBuilder{}

   //add openids  (if you want to update other people's card otherwise you can ignore it)
   openids := make([]string,1)
   openids[0] = "another user's openid"
   builder.OpenIDs = openids

   //add config
   config := protocol.ConfigForm{
   MinVersion:     protocol.VersionForm{},
   WideScreenMode: true,
   EnableForward: false,
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
