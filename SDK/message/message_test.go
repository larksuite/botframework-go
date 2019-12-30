// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package message_test

import (
	"context"
	"testing"

	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

func TestSendTextMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}
	text := "sendTextMessage test(~_~)"
	resp, err := message.SendTextMessage(c, tenantKey, appConf.AppID, user, "", text)
	if err != nil {
		t.Errorf("SendTextMessage: failed err[%v]", err)
	} else {
		t.Logf("SendTextMessage: succ messageID[%s]", resp.Data.MessageID)
	}
}

func TestSendImageMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()
	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}
	//choose url,path or imageKey
	//by url
	url := "https://s0.pstatp.com/ee/lark-open/web/static/apply.226f11cb.png"
	resp, err := message.SendImageMessage(c, tenantKey, appConf.AppID, user, "", url, "", "")
	if err != nil {
		t.Errorf("SendImageMessage: failed err[%v]", err)
	} else {
		t.Logf("SendImageMessage: succ messageID[%s]", resp.Data.MessageID)
	}
	//by path
	path := "../../demo/source/lark0.jpg"
	resp, err = message.SendImageMessage(c, tenantKey, appConf.AppID, user, "", "", path, "")
	if err != nil {
		t.Errorf("SendImageMessage: failed err[%v]", err)
	} else {
		t.Logf("SendImageMessage: succ messageID[%s]", resp.Data.MessageID)
	}
	//by imageKey
	resp, err = message.SendImageMessage(c, tenantKey, appConf.AppID, user, "", "", "", imageKey)
	if err != nil {
		t.Errorf("SendImageMessage: failed err[%v]", err)
	} else {
		t.Logf("SendImageMessage: succ messageID[%s]", resp.Data.MessageID)
	}
}

func TestSendRichTextMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	postForm := make(map[protocol.Language]*protocol.RichTextForm)

	// en-us
	titleUS := "this is a title"
	contentUS := message.NewRichTextContent()

	// first line
	contentUS.AddElementBlock(
		message.NewTextTag("first line: ", true, 1),
		message.NewATag("hyperlinks ", true, "https://www.feishu.cn"),
		message.NewAtTag("username", userID),
	)

	// second line
	contentUS.AddElementBlock(
		message.NewTextTag("second line: ", true, 1),
		message.NewTextTag("text test", true, 1),
	)

	postForm[protocol.EnUS] = message.NewRichTextForm(&titleUS, contentUS)

	// zh-cn
	titleCN := "这是一个标题"
	contentCN := message.NewRichTextContent()

	// first line
	contentCN.AddElementBlock(
		message.NewTextTag("第一行: ", true, 1),
		message.NewATag("超链接 ", true, "https://www.feishu.cn"),
		message.NewAtTag("username", userID),
	)

	// second line
	contentCN.AddElementBlock(
		message.NewTextTag("第二行: ", true, 1),
		message.NewTextTag("文本测试", true, 1),
	)

	postForm[protocol.ZhCN] = message.NewRichTextForm(&titleCN, contentCN)

	resp, err := message.SendRichTextMessage(c, tenantKey, appConf.AppID, user, "", postForm)
	if err != nil {
		t.Errorf("SendRichTextMessage: failed err[%v]", err)
	} else {
		t.Logf("SendRichTextMessage: succ messageID[%s]", resp.Data.MessageID)
	}
}

func TestSendShareChatMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendShareChatMessage(c, tenantKey, appConf.AppID, user, "", chatID)
	if err != nil {
		t.Errorf("SendShareChatMessage: failed err[%v]", err)
	} else {
		t.Logf("SendShareChatMessage: succ messageID[%s]", resp.Data.MessageID)
	}
}

func TestSendCardMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	//card builder
	builder := &message.CardBuilder{}

	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)

	//add header
	content := "Please choose color"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")

	//add hr
	builder.AddHRBlock()

	//add block
	builder.AddDIVBlock(nil, []protocol.FieldForm{
		*message.NewField(false, message.NewMDText("**Async**", nil, nil, nil)),
	}, nil)

	//add divBlock
	builder.AddDIVBlock(nil, []protocol.FieldForm{
		*message.NewField(false, message.NewMDText("**Sync**", nil, nil, nil)),
	}, nil)

	//add actionBlock
	payload1 := make(map[string]string, 0)
	payload1["color"] = "red"
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewButton(message.NewMDText("red", nil, nil, nil),
			nil, nil, payload1, protocol.PRIMARY, nil, "asyncButton"),
	})

	//add jumpBlock
	url := "https://www.google.com"
	ext := message.NewJumpButton(message.NewMDText("jump to google", nil, nil, nil), &url, nil, protocol.DEFAULT)
	builder.AddDIVBlock(message.NewMDText("", nil, nil, nil), nil, ext)

	//add imageBlock
	builder.AddImageBlock(
		message.NewMDText("", nil, nil, nil),
		*message.NewMDText("", nil, nil, nil),
		imageKey)

	//generate card
	card, err := builder.BuildForm()

	resp, err := message.SendCardMessage(c, tenantKey, appConf.AppID, user, "", *card, true)
	if err != nil {
		t.Errorf("SendCardMessage: failed err[%v]", err)
	} else {
		t.Logf("SendCardMessage: succ messageID[%s]", resp.Data.MessageID)
	}
}

//need to get access
func TestBatchSendTextMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()

	info := &protocol.BatchBaseInfo{
		DepartmentIDs: []string{},
		OpenIDs:       []string{openID},
		UserIDs:       []string{},
	}
	text := "sendTextMessage test(~_~)"
	resp, err := message.SendTextMessageBatch(c, tenantKey, appConf.AppID, info, "", text)
	if err != nil {
		t.Errorf("SendTextMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendTextMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}
}

//need to get access
func TestBatchSendImageMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()

	info := &protocol.BatchBaseInfo{
		DepartmentIDs: []string{},
		OpenIDs:       []string{openID},
		UserIDs:       []string{},
	}

	//choose url,path or imageKey
	//by url
	url := "https://s0.pstatp.com/ee/lark-open/web/static/apply.226f11cb.png"
	resp, err := message.SendImageMessageBatch(c, tenantKey, appConf.AppID, info, "", url, "", "")
	if err != nil {
		t.Errorf("SendImageMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendImageMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}

	//by path
	path := "../../demo/source/lark0.jpg"
	resp, err = message.SendImageMessageBatch(c, tenantKey, appConf.AppID, info, "", "", path, "")
	if err != nil {
		t.Errorf("SendImageMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendImageMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}

	//by imageKey
	resp, err = message.SendImageMessageBatch(c, tenantKey, appConf.AppID, info, "", "", "", imageKey)
	if err != nil {
		t.Errorf("SendImageMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendImageMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}
}

//need to get access
func TestBatchSendRichTextMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()

	info := &protocol.BatchBaseInfo{
		DepartmentIDs: []string{},
		OpenIDs:       []string{openID},
		UserIDs:       []string{},
	}

	postForm := make(map[protocol.Language]*protocol.RichTextForm)

	// en-us
	titleUS := "this is a title"
	contentUS := message.NewRichTextContent()

	// first line
	contentUS.AddElementBlock(
		message.NewTextTag("first line: ", true, 1),
		message.NewATag("hyperlinks ", true, "https://www.feishu.cn"),
		message.NewAtTag("username", userID),
	)

	// second line
	contentUS.AddElementBlock(
		message.NewTextTag("second line: ", true, 1),
		message.NewTextTag("text test", true, 1),
	)

	postForm[protocol.EnUS] = message.NewRichTextForm(&titleUS, contentUS)

	// zh-cn
	titleCN := "这是一个标题"
	contentCN := message.NewRichTextContent()

	// first line
	contentCN.AddElementBlock(
		message.NewTextTag("第一行: ", true, 1),
		message.NewATag("超链接 ", true, "https://www.feishu.cn"),
		message.NewAtTag("username", userID),
	)

	// second line
	contentCN.AddElementBlock(
		message.NewTextTag("第二行: ", true, 1),
		message.NewTextTag("文本测试", true, 1),
	)

	postForm[protocol.ZhCN] = message.NewRichTextForm(&titleCN, contentCN)

	resp, err := message.SendRichTextMessageBatch(c, tenantKey, appConf.AppID, info, "", postForm)
	if err != nil {
		t.Errorf("SendRichTextMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendRichTextMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}
}

//need to get access
func TestBatchSendShareChatMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()
	info := &protocol.BatchBaseInfo{
		DepartmentIDs: []string{},
		OpenIDs:       []string{openID},
		UserIDs:       []string{},
	}
	resp, err := message.SendShareChatMessageBatch(c, tenantKey, appConf.AppID, info, "", chatID)
	if err != nil {
		t.Errorf("SendShareChatMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendShareChatMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}
}

//need to get access
func TestBatchSendCardMessage(t *testing.T) {
	c := context.Background()
	InitTestParams()

	info := &protocol.BatchBaseInfo{
		DepartmentIDs: []string{},
		OpenIDs:       []string{openID},
		UserIDs:       []string{},
	}

	url := "https://www.google.com"
	content := "card message test"
	card := protocol.CardForm{
		Config: &protocol.ConfigForm{MinVersion: protocol.VersionForm{Version: "1.0"},
			Debug:          true,
			WideScreenMode: true,
		},
		CardLink: &protocol.URLForm{
			Url: &url,
		},
		Header: &protocol.CardHeaderForm{
			Title: protocol.TextForm{
				Tag:     protocol.PLAIN_TEXT_E,
				Content: &content,
			},
			Template: "",
		},
		Elements: []interface{}{},
	}

	resp, err := message.SendCardMessageBatch(c, tenantKey, appConf.AppID, info, "", card, true)
	if err != nil {
		t.Errorf("SendCardMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendCardMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}
}
