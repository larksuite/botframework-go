# 消息卡片的构造  
卡片主要分为三个部分： config 、 header 和 elements 。  
**config**:  
config 主要存储一系列卡片配置。config 模块中只有一个 bool 类型的字段 wide_screen_mode 。表示是否根据屏幕宽度动态调整消息卡片的宽度。默认为 false 值。  
**header**:  
header 主要配置消息卡片的标题内容。只有一个字段 title ，title 是一个 text 对象，且仅支持 "plain_text" 。    
**element**:  
element 为一个个模块。卡片的内容正是由一个个模块堆砌起来的。模块又分为内容模块、分割线模块、图片模块、交互模块、备注模块等五种。     
  
## 如何构造并发送一张卡片  
这里通过示例代码详细介绍如何构造一张卡片。    
1. 声明一个结构体用来存储卡片的信息  
```go
builder := &message.CardBuilder{}
```
  
2. 添加和配置 config    
```go
//add config
type ConfigForm struct {
	MinVersion     VersionForm `json:"min_version,omitempty" validate:"omitempty"`
	Debug          bool        `json:"debug,omitempty" validate:"omitempty"`
	WideScreenMode bool        `json:"wide_screen_mode,omitempty" validate:"omitempty"`
	EnableForward bool        `json:"enable_forward,omitempty" validate:"omitempty"`
}
builder.SetConfig(config)
```
  
3. 添加 header  
```go
//add header
content := "this is a card"
line := 1
title := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &content,
   Lines:   &line,
}
builder.AddHeader(title, "")
```
  
4. 添加一个简单的模块   
```go
//add button
button1 := make(map[string]string, 0)
button1["key"] = "buttonValue"
builder.AddActionBlock([]protocol.ActionElement{
   message.NewButton(message.NewMDText("callback button", nil, nil, nil), nil, nil, button1, protocol.DANGER, nil, "clickbutton"),
})
```
  
5. 将该结构体转换为卡片类型，并通过 message.SendCardMessage 函数发送卡片   
```go
card, err := builder.BuildForm()
if err != nil {
   return fmt.Errorf("card build failed error[%v]", err)
}
_, err = message.SendCardMessage(ctx, tenantKey, appID, user, "", *card, false)
if err != nil {
   return fmt.Errorf("send message failed error[%v]", err)
}
```
  
## 模块说明  
### 内容模块  
内容模块是用途最灵活的模块，可以在其中组合各类用于展示的内容。    
模块标签为 “div” ，可以单独通过 text 和 field 来展示文本内容，也可以配合一个 image 元素或一个 button ， overflow ， selectMenu ， pickerDate 等互动元素增加内容的丰富性。    
```go
//设置文本的超链接  如果不需要可以不设置
urlhref := make(map[string]protocol.URLForm, 0)
uHost := "https://www.larksuite.com/"
urlhref["link"] = protocol.URLForm{Url: &uHost}
//增加内容模块的标题，内容，附加元素（这里没有添加）
builder.AddDIVBlock(message.NewMDText("[lark]($link)", nil, nil, urlhref), []protocol.FieldForm{
   *message.NewField(false, message.NewMDText("**boldText**", nil, nil, nil)),
   *message.NewField(false, message.NewMDText("text", nil, nil, nil)),
}, nil)
```
  
### 分割线模块  
分割线模块在模块之间添加一条分割线。  
```go
builder.AddHRBlock()
```
  
### 图片模块
图片模块分为三个部分: title，alt和imageky。    
tilte中存放图片的标题。alt中存放鼠标指针悬停在图片上时显示的图片说明。  
```go
//增加图片模块
ImageContent := "Description when your mouse is over the picture"
ImageTag := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &ImageContent,
   Lines:   nil,
}

//获取imageKey，可以在 docs/zh/send_message.md/机器人发送图片消息 查看imagekey的详细说明
imageURL := "https://s0.pstatp.com/ee/lark-open/web/static/apply.226f11cb.png"
imageKey, err := message.GetImageKey(ctx, tenantKey, appID, imageURL, "")
if err != nil {
   return fmt.Errorf("get imageKey failed[%v]", err)
}
builder.AddImageBlock(message.NewMDText("image title", nil, nil, nil), ImageTag, imageKey)
```
  
### 备注模块  
备注模块用来展示用于辅助说明或备注的次要信息，支持小尺寸的图片和文本。  
```go
//增加备注模块
noteTextContent := "备注文本内容"
notetext := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &noteTextContent,
   Lines:   nil,
}
noteImageContent := "备注图片内容"
noteImage := protocol.ImageForm{
   Tag:      protocol.IMG_BLOCK,
   ImageKey: imageKey, //这里填写正确的imageKey
   ALT: protocol.TextForm{
      Tag:     protocol.PLAIN_TEXT_E,
      Content: &noteImageContent,
      Lines:   nil,
   },
}
builder.AddNoteBlock([]protocol.BaseElement{
   &notetext,
   &noteImage,
})
```
  
#### 交互模块  
交互模块用于承载交互元素。    
卡片提供了 4 类交互元素，通过 elements 字段添加交互元素，实现交互功能。    

回传方式：每个交互元素或交互选项都提供了 value 字段，用户点击交互元素或选项之后，业务方会收到对应的 value 字段值，以此决定后续操作。业务方可知用户行为。    

跳转方式：button 和 overflow 的选项配置跳转链接，用户点击跳转相应链接，不回传 value 字段值到业务方。业务方不可知用户行为。    

#### button  
button属于交互元素的一种，可用于内容块的extra字段和交互块的elements字段。    
```go
button1 := make(map[string]string, 0)
button1["key"] = "value"
url := "https://www.feishu.cn"
builder.AddActionBlock([]protocol.ActionElement{
message.NewButton(message.NewMDText("button",nil,nil,nil),&url,nil,button1,protocol.PRIMARY,nil,"nil"),
})
```
  
#### selectMenu  
selectMenu属于交互元素的一种，提供选项菜单的功能，可用于内容块的extra字段和交互块的elements字段。    
选项菜单  
```go
//SelectStaticMenu
optiontextcontent1 := "选项一"
optiontext1 := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &optiontextcontent1,
   Lines:   nil,
}
optiontextcontent2 := "选项二"
optiontext2 := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &optiontextcontent2,
   Lines:   nil,
}
value1 := "value1"
value2 := "value2"
builder.AddActionBlock([]protocol.ActionElement{
   message.NewSelectStaticMenu(message.NewMDText("SelectStaticMenu", nil, nil, nil), nil, []protocol.OptionForm{
      message.NewOption(optiontext1, value1),
      message.NewOption(optiontext2, value2),
   }, &value1, nil, "testcard"), //option1 is default option
},
)
```
    
选人菜单  
```go
//SelectPersonMenu
builder.AddActionBlock([]protocol.ActionElement{
   message.NewSelectPersonMenu(message.NewMDText("SelectPersonMenu", nil, nil, nil), nil, []protocol.OptionForm{}, nil, nil, "testcard"),
},
)
```
  
#### overflow  
提供折叠的按钮型菜单。    
overflow 属于交互元素的一种，可用于内容块的 extra 字段和交互块的 elements 字段。  
通过 options 字段配置选项，可用于多个按扭的折叠隐藏功能。    
```go
//add overflow
overFlowcontent1 := "选项一"
overFlowtext1 := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &overFlowcontent1,
   Lines:   nil,
}
overFlowcontent2 := "选项二"
overFlowtext2 := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &overFlowcontent2,
   Lines:   nil,
}
builder.AddActionBlock([]protocol.ActionElement{
   message.NewOverflowMenu(nil, []protocol.OptionForm{
      message.NewOption(overFlowtext1, value1),
      message.NewOption(overFlowtext2, value2),
   }, nil, "testcard"),
},
)
```
  
##### PickerDate  
提供日期选择的功能。    
PickerDate 属于交互元素的一种，可用于内容块的extra字段和交互块的elements字段。    
```go
//add PickerDate
timePicker := time.Now().Format("2006-01-02")
builder.AddActionBlock([]protocol.ActionElement{
   message.NewPickerDate(message.NewMDText("PickerDate", nil, nil, nil), nil, nil, &timePicker, "testcard"),
},
)
```
  
##### 延迟更新卡片  
业务方使用交互返回的 token 凭证，在30分钟内最多更新两次卡片。  
除正常卡片内容外，还有 open_ids 字段控制更新指定用户的消息卡片，字段值为用户 openId 数组。  
示例代码如下：  
```go
// use message.UpdateCard function to update card
func UpdateCard(token, tenantkey, appid string, openid []string) error {
    //token可以在交互模块回调时获得

	//build a new card
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

	builder.AddDIVBlock(nil, []protocol.FieldForm{
		*message.NewField(false, message.NewMDText("updatecard by message.UpdateCard", nil, nil, nil)),
	}, nil)

	card, err := builder.BuildForm()
	if err != nil {
		return fmt.Errorf("build card failed error[%v]", err)
	}

	// card.OpenIDs can't be nil.
	card.OpenIDs = openid

	_, err = message.UpdateCard(context.TODO(), tenantkey, appid, token, *card)
	if err != nil {
		return fmt.Errorf("card update failed error[%v]", err)
	}
	return nil
}
```

### 发送卡片 demo  
[demo code](../../demo/send_card/send_card.go)   
该示例代码实现了卡片内增加内容模块、分割线、图片模块、交互模块、备注模块等，可供开发者参考。  
