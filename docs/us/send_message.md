# bot sends all kinds of messages  
You can use functions in SDK to send all kinds of messages easily.  

## text  
Example code:  
```go
//send text
func sendTextMessage(chatID, tenantKey, appID string) error {
    // Necessary steps. Necessary params. ID and type correspond one by one
    user := &protocol.UserInfo{
        ID:   chatID,
        Type: protocol.UserTypeChatID,
    }

    // Necessary steps. Necessary params.
    ctx := context.TODO()

    //Penultimate param means rootid，if you want to reply a message,rootid is this message's id.You can use "" in general.
    _, err := message.SendTextMessage(ctx, tenantKey, appID, user, "", "在飞书，享高效")
    if err != nil {
        return fmt.Errorf("send text failed[%v]", err)
    }

    return nil
}
```
  
## image  
There are three ways to send pictures: imageKey, path and url.（priority: imageKey>path>url）  
Example code:  
```go
//send image
func sendImageMessage(chatID, tenantKey, appID string) error {
   //Necessary steps. Necessary params. ID and type correspond one by one
    user := &protocol.UserInfo{
        ID:   chatID,
        Type: protocol.UserTypeChatID,
    }
    //Necessary steps. Necessary params.
    ctx := context.TODO()

    //send image(use imageurl)
    url := "https://is3-ssl.mzstatic.com/image/thumb/Purple113/v4/ed/ee/c0/edeec03e-d111-ac8d-3441-409acd11dbea/source/512x512bb.jpg"
    _, err := message.SendImageMessage(ctx, tenantKey, appID, user, "", url, "", "")
    if err != nil {
        return fmt.Errorf("send image failed[%v]", err)
    }
    return nil
}
```
  
If you want to use imageKey to send an image, you should get imageKey first.  
Example code:  
```go
    user := &protocol.UserInfo{
        ID:   chatID,
        Type: protocol.UserTypeChatID,
    }
    ctx := context.TODO()

    url := "https://is3-ssl.mzstatic.com/image/thumb/Purple113/v4/ed/ee/c0/edeec03e-d111-ac8d-3441-409acd11dbea/source/512x512bb.jpg"
    imagekey,err := message.GetImageKey(ctx, tenantKey, appID, url, "" )
    if err != nil{
        return fmt.Errorf("get imageKey failed[%v]", err)
    }
    _, err = message.SendImageMessage(ctx, tenantKey, appID, user, "", "", "", imagekey)
    if err != nil {
        return fmt.Errorf("send image failed[%v]", err)
    }
```

## rich text  
Example code:    
```go
//send rich text
func sendRichTextMessage(chatID, tenantKey, appID, userID string) error {
    user := &protocol.UserInfo{
        ID:   chatID,
        Type: protocol.UserTypeChatID,
    }
    ctx := context.TODO()

    //add content of richtext
    //zh-cn
    titleCN := "这是一个标题"
    contentCN := message.NewRichTextContent()
    // first line
    contentCN.AddElementBlock(
        message.NewTextTag("第一行 :", true, 1),
        message.NewATag("超链接", true, "https://www.feishu.cn"),
        message.NewAtTag("用户名", userID),
    )
    // second line
    contentCN.AddElementBlock(
        message.NewTextTag("第二行 :", true, 1),
        message.NewTextTag("文本测试", true, 1),
    )

    //en-us
    titleUS := "this is a title"
    contentUS := message.NewRichTextContent()
    // first line
    contentUS.AddElementBlock(
        message.NewTextTag("first line :", true, 1),
        message.NewAtTag("username", userID),
    )
    // second line
    contentUS.AddElementBlock(
        message.NewTextTag("second line :", true, 1),
        message.NewTextTag("text test", true, 1),
    )

    postForm := make(map[protocol.Language]*protocol.RichTextForm)
    postForm[protocol.ZhCN] = message.NewRichTextForm(&titleCN, contentCN)
    postForm[protocol.EnUS] = message.NewRichTextForm(&titleUS, contentUS)

    //send rich text
    _, err := message.SendRichTextMessage(ctx, tenantKey, appID, user, "", postForm)
    if err != nil {
        return fmt.Errorf("send rich text failed[%v]", err)
    }
    return nil
}
```
  
## group card  
Example code:  
```go
//send group shared card
func sendShareChatMessage(openID, tenantKey, appID, shareChatID string) error {
    user := &protocol.UserInfo{
        ID:   openID,
        Type: protocol.UserTypeOpenID,
    }

    ctx := context.TODO()

    //send group shared card(last param means group chat id, and this message will be sent to this user by openid)
    _, err := message.SendShareChatMessage(ctx, tenantKey, appID, user, "", shareChatID)
    if err != nil {
        return fmt.Errorf("send group card failed[%v]", err)
    }
    return nil
}
```
  
### send message demo  
[demo code](../../demo/send_message/send_message.go)   
