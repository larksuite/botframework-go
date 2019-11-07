# card
**config**:  
Config mainly stores a series of card configurations. There is only one bool type field wide_screen_mode in the config module. This field determines whether to dynamically adjust the width of the message card based on the screen width. The default value is false.  
  
**header**:    
The header mainly configures the header content of the message card. There is only one field, title, a plain_text object.  

**element**:    
The contents of the card are piled up by modules. The module is divided into five modules: content module, split line module, image module, interaction module and note module.  

## how to build and send a card  
Example code:  
1. Declare an struct to store card information.    
```go
builder := &message.CardBuilder{}
```

2. Add config    
```go
//add config
config := protocol.ConfigForm{
   MinVersion:     protocol.VersionForm{},
   WideScreenMode: true,
}
builder.SetConfig(config)
```
  
3. Add header  
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

4. Add a simple block  
```go
//add button
button1 := make(map[string]string, 0)
button1["key"] = "buttonValue"
builder.AddActionBlock([]protocol.ActionElement{
   message.NewButton(message.NewMDText("callback button", nil, nil, nil), nil, nil, button1, protocol.DANGER, nil, "clickbutton"),
})
``` 
  
5. Declare an object of card  
```go
card, err := builder.BuildForm()
if err != nil {
   return fmt.Errorf("card build failed error[%v]", err)
}
```

6. Send a card with true params  
```go
_, err := message.SendCardMessage(ctx, tenantKey, appID, user, "", *card, false)
if err != nil {
   return fmt.Errorf("send message failed error[%v]", err)
}
```
  
## how to add other modules  

### content module  
Content module is the most flexible module, in which you can combine all kinds of content for display.    
The module label is "div", which can display the text content through text and field alone, or add the richness of content with an image element or other interactive elements.    
Example code:  
```go
//set href
urlhref := make(map[string]protocol.URLForm, 0)
uHost := "https://www.larksuite.com/"
urlhref["link"] = protocol.URLForm{Url: &uHost}

//add content block,suck as title、content and extra（null here）.
builder.AddDIVBlock(message.NewMDText("[lark]($link)", nil, nil, urlhref), []protocol.FieldForm{
   *message.NewField(false, message.NewMDText("**boldText**", nil, nil, nil)),
   *message.NewField(false, message.NewMDText("text", nil, nil, nil)),
}, nil)
```
  
### split line module  
The divider module can add a divider between the modules.  
Example code:  
```go
builder.AddHRBlock()
```

### image module  
The module is divided into three parts: title, alt and imageky.        
Title means the title of an image. ALT stores the picture description displayed when the mouse pointer hovers over the picture.     
For an explanation of the imagekey field, refer to bot sends all kinds of messages-image.  
Example code:  
```go
//add image block
ImageContent := "Description when your mouse is over the picture"
ImageTag := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &ImageContent,
   Lines:   nil,
}

//get imageKey.You can get detailed introduction about imageKey from send_message.md/image
imageURL := "https://s0.pstatp.com/ee/lark-open/web/static/apply.226f11cb.png"
imageKey, err := message.GetImageKey(ctx, tenantKey, appID, imageURL, "")
if err != nil {
   return fmt.Errorf("get imageKey failed[%v]", err)
}
builder.AddImageBlock(message.NewMDText("image title", nil, nil, nil), ImageTag, imageKey)
```
  
### note module  
Note module is used to display the secondary information used for auxiliary description or note, supporting small size pictures and texts.  
Example code:  
```go
//add note
noteTextContent := "noteTextContent"
notetext := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &noteTextContent,
   Lines:   nil,
}
noteImageContent := "noteImageContent"
noteImage := protocol.ImageForm{
   Tag:      protocol.IMG_BLOCK,
   ImageKey: imageKey, //you should use true imageKey
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

### interaction module  
The card provides four types of interactive elements, which are added through the elements field to realize the interactive function.    
  
Return method: each interaction element or option provides a value field. After the user clicks the interaction element or option, the business side will receive the corresponding value field value to determine the subsequent operation. The business can know the user's behavior.  

Jump method: the option of button and overflow configures the jump link. The user clicks the corresponding link and does not return the value field value to the business party. The business can not know the user behavior.  

#### button  
Example code:  
```go
//add button
uHost := "https://www.larksuite.com/"
button1 := make(map[string]string, 0)
button1["key"] = "value"
builder.AddActionBlock([]protocol.ActionElement{
   message.NewButton(message.NewMDText("lark button", nil, nil, nil), &uHost, nil, button1, protocol.PRIMARY, nil, "testcard"),
   message.NewButton(message.NewMDText("simple button", nil, nil, nil), nil, nil, button1, protocol.DANGER, nil, "testcard"),
})//The penultimate param means confirm.
```

#### selectMenu  
Example code:  

SelectStaticMenu   
```go
//SelectStaticMenu
optiontextcontent1 := "option1"
optiontext1 := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &optiontextcontent1,
   Lines:   nil,
}
optiontextcontent2 := "option2"
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

SelectPersonMenu  
```go
//SelectPersonMenu
builder.AddActionBlock([]protocol.ActionElement{
   message.NewSelectPersonMenu(message.NewMDText("SelectPersonMenu", nil, nil, nil), nil, []protocol.OptionForm{}, nil, nil, "testcard"),
},
)
```

  
#### overflow  
Collapsed button menu.  
Example code:  
```go
//add overflow
overFlowcontent1 := "option1"
overFlowtext1 := protocol.TextForm{
   Tag:     protocol.PLAIN_TEXT_E,
   Content: &overFlowcontent1,
   Lines:   nil,
}
overFlowcontent2 := "option2"
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
  
#### PickerDate  
Example code:  
```go
//add PickerDate
timePicker := time.Now().Format("2006-01-02")
builder.AddActionBlock([]protocol.ActionElement{
   message.NewPickerDate(message.NewMDText("PickerDate", nil, nil, nil), nil, nil, &timePicker, "testcard"),
},
)
```

##### Delayed update card  
You can use token to update the card at most twice in 30 minutes.   
In addition to the normal card content, the OpenIDs field controls the update of the message card of the specified user.     
The field value is the user openid array.       
Example code:    
```go
// use message.UpdateCard function to update card
func UpdateCard(token,tenantkey,appid,openid string)(error){
    //you can get token by card callback

	//build a new card
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

	builder.AddDIVBlock(nil, []protocol.FieldForm{
		*message.NewField(false, message.NewMDText("updatecard by message.UpdateCard", nil, nil, nil)),
	}, nil)

	card, err := builder.BuildForm()
	if err != nil {
		return fmt.Errorf("build card failed error[%v]", err)
	}

	// card.OpenIDs can't be nil.
	card.OpenIDs=[]string{openid}

	_,err = message.UpdateCard(context.TODO(),tenantkey,appid,token,*card)
	if err != nil{
		return fmt.Errorf("card update failed error[%v]", err)
	}
	return nil
}
```
  

### send card demo  
[demo code](../../demo/send_card/send_card.go)
