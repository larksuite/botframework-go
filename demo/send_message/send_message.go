package main

import (
	"context"
	"fmt"

	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
	"github.com/larksuite/botframework-go/demo/sdk_init"
)

func main() {
	//init log
	common.InitLogger(common.NewCommonLogger(), common.DefaultOption())
	defer common.FlushLogger()

	//param
	ctx := context.TODO()

	//Necessary step: init app configuration
	err := sdk_init.InitInfo()
	if err != nil {
		common.Logger(ctx).Errorf("init app config failed[%v]", err)
		return
	}

	//params
	chatID := ""    //p2p or group chat ID
	tenantKey := "" //tenantKey of your company
	appID := ""     //APP ID

	//send text
	sendTextMessage(chatID, tenantKey, appID)

	//send image
	sendImageMessage(chatID, tenantKey, appID)

	//send rich text
	//param
	userID := "" // @this user
	sendRichTextMessage(chatID, tenantKey, appID, userID)

	//send group card
	//params
	openID := ""     //User_ID in this app.Group card will be sent to this user who has this openID.
	sharChatID := "" //group's chatid

	sendShareChatMessage(openID, tenantKey, appID, sharChatID)

	//send card
	//you can find the example in demo/send_card/send_card.go
}

//send text
func sendTextMessage(chatID, tenantKey, appID string) error {
	user := &protocol.UserInfo{
		ID:   chatID,
		Type: protocol.UserTypeChatID,
	}

	ctx := context.TODO()

	_, err := message.SendTextMessage(ctx, tenantKey, appID, user, "", "Always One, Always Agile")
	if err != nil {
		common.Logger(ctx).Errorf("send text failed[%v]", err)
		return fmt.Errorf("send text failed[%v]", err)
	}

	return nil
}

//send image
func sendImageMessage(chatID, tenantKey, appID string) error {
	user := &protocol.UserInfo{
		ID:   chatID,
		Type: protocol.UserTypeChatID,
	}

	ctx := context.TODO()

	//send image(use imageurl)
	imageURL := "https://is3-ssl.mzstatic.com/image/thumb/Purple113/v4/ed/ee/c0/edeec03e-d111-ac8d-3441-409acd11dbea/source/512x512bb.jpg"
	_, err := message.SendImageMessage(ctx, tenantKey, appID, user, "", imageURL, "", "")
	if err != nil {
		common.Logger(ctx).Errorf("send image failed[%v]", err)
		return fmt.Errorf("send image failed[%v]", err)
	}
	return nil
}

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
		common.Logger(ctx).Errorf("send rich text failed[%v]", err)
		return fmt.Errorf("send rich text failed[%v]", err)
	}
	return nil
}

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
		common.Logger(ctx).Errorf("send group card failed[%v]", err)
		return fmt.Errorf("send group card failed[%v]", err)
	}
	return nil
}
