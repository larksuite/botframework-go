# 机器人发送各类消息  
可以非常简单的通过调用 SDK 中的一系列函数来发送各类消息。  
  
## 机器人发送文本信息  
调用 SDK 中的 message.SendTextMessage 函数即可发送文本消息。      
示例代码如下:  
 ```go
func sendTextMessage(chatID, tenantKey, appID string) error {
   user := &protocol.UserInfo{
      ID:   chatID,
      Type: protocol.UserTypeChatID,
   }
   ctx := context.TODO()

   _, err := message.SendTextMessage(ctx, tenantKey, appID, user, "", "在飞书，享高效")
   if err != nil {
      return fmt.Errorf("send text failed[%v]", err)
   }

   return nil
}
```
  
## 机器人发送图片信息  
调用 SDK 中的 message.SendImageMessage 函数即可发送图片消息。    
该函数的最后三个参数分别为图片 url 、图片 path 和图片 imageKey 。
这三个参数填写任意一个即可，如果填写多个，则会按照 imagekey -> path -> url 的优先级使用。
底层接口实际需要 imagekey ，当使用 path 或 url 时，函数内部会先调用上传图片接口换取 imagekey 。
示例代码如下：  
```go
//发送图片消息函数
func sendImageMessage(chatID, tenantKey, appID string) error {
    user := &protocol.UserInfo{
        ID:   chatID,
        Type: protocol.UserTypeChatID,
    }
   ctx := context.TODO()

    //发送图片消息，以url为例
    url := "https://s0.pstatp.com/ee/lark-open/web/static/apply.226f11cb.png"
    _, err := message.SendImageMessage(ctx, tenantKey, appID, user, "", url, "", "")
    if err != nil {
        return fmt.Errorf("send image failed[%v]", err)
    }
    return nil
}
```
  
如果想通过 imagekey 发送图片消息，需要先获取图片的 imagkey 。（获取 imagkey 的函数仍然需要 url 或 path 其中一项）    
先上传图片，得到 imagekey ，以后就可以通过该 imagekey 发送图片了。    
实际上 message.SendImageMessage 函数内部也是通过该方式先获取 imagekey 再通过 imagekey 发送图片的。  
 ```go
    user := &protocol.UserInfo{
        ID:   chatID,
        Type: protocol.UserTypeChatID,
    }
    ctx := context.TODO()

    url := "https://s0.pstatp.com/ee/lark-open/web/static/apply.226f11cb.png"
    imagekey,err := message.GetImageKey(ctx, tenantKey, appID, url, "" )
    if err != nil{
        return fmt.Errorf("get imageKey failed[%v]", err)
    }

    _, err = message.SendImageMessage(ctx, tenantKey, appID, user, "", "", "", imagekey)
    if err != nil {
        return fmt.Errorf("send image failed[%v]", err)
    }
```

## 机器人发送富文本信息  
调用 SDK 中的 message.SendRichTextMessage 函数即可发送富文本消息。    
示例代码如下：    
```go
//发送富文本消息函数
func sendRichTextMessage(chatID, tenantKey, appID, userID string) error {
    user := &protocol.UserInfo{
        ID:   chatID,
        Type: protocol.UserTypeChatID,
    }
    ctx := context.TODO()

    //构建富文本内容

    //i18n zh-CN
    //标题内容自定义，将作为必要参数传入
    titleCN := "这是一个标题"
    //设置富文本content，将作为必要参数传入
    contentCN := message.NewRichTextContent()
    //设置富文本具体内容，如果不设置则没有内容，用户自定义即可
    contentCN.AddElementBlock(
        message.NewTextTag("第一行 :", true, 1),
        message.NewATag("超链接", true, "https://www.feishu.cn"),
        message.NewAtTag("用户名", userID),
    )
    contentCN.AddElementBlock(
        message.NewTextTag("第二行 :", true, 1),
        message.NewTextTag("文本测试", true, 1),
    )

    //i18n en-US
    titleUS := "this is a title"
    contentUS := message.NewRichTextContent()
    contentUS.AddElementBlock(
        message.NewTextTag("first line :", true, 1),
        message.NewATag("href", true, "https://www.feishu.cn"),
        message.NewAtTag("username", userID),
    )
    contentUS.AddElementBlock(
        message.NewTextTag("second line :", true, 1),
        message.NewTextTag("text test", true, 1),
    )
    postForm := make(map[protocol.Language]*protocol.RichTextForm)
    postForm[protocol.ZhCN] = message.NewRichTextForm(&titleCN, contentCN)
    postForm[protocol.EnUS] = message.NewRichTextForm(&titleUS, contentUS)
    
    //发送富文本消息
    _, err := message.SendRichTextMessage(ctx, tenantKey, appID, user, "", postForm)
    if err != nil {
        return fmt.Errorf("send rich text failed[%v]", err)
    }
    return nil
}
```
  
## 机器人发送群名片信息  
调用 SDK 中的 message.SendRichTextMessage 函数即可发送群名片消息。    
示例代码如下：     
```go
//发送群名片消息函数
func sendShareChatMessage(openID, tenantKey, appID, shareChatID string) error {
   user := &protocol.UserInfo{
      ID:   openID,
      Type: protocol.UserTypeOpenID,
   }

   ctx := context.TODO()

   _, err := message.SendShareChatMessage(ctx, tenantKey, appID, user, "", shareChatID)
   if err != nil {
      return fmt.Errorf("send group card failed[%v]", err)
   }
   return nil
}
```
  
## 发送卡片消息  
调用 SDK 中的 message.SendCardMessage 函数即可发送卡片消息。    
[demo code](../../demo/send_card/send_card.go)
该示例代码实现了卡片内增加内容模块、分割线、图片模块、交互模块、备注模块等，可供开发者参考。  
