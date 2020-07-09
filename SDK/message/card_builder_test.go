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

func TestCardAsyncButton(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	builder := &message.CardBuilder{}
	//add header
	content := "Please choose color"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")
	builder.AddHRBlock()
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)
	//add block
	builder.AddDIVBlock(nil, []protocol.FieldForm{
		*message.NewField(false, message.NewMDText("**Async**", nil, nil, nil)),
	}, nil)
	payload1 := make(map[string]string, 0)
	payload1["color"] = "red"
	payload2 := make(map[string]string, 0)
	payload2["color"] = "black"
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewButton(message.NewMDText("red", nil, nil, nil),
			nil, nil, payload1, protocol.PRIMARY, nil, "asyncButton"),
		message.NewButton(message.NewMDText("black", nil, nil, nil),
			nil, nil, payload2, protocol.PRIMARY, nil, "asyncButton"),
	})
	card, err := builder.BuildForm()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, false)
	if err != nil {
		t.Errorf("TestAsyncButton failed:%v\n", err)
	} else {
		t.Logf("code:%d ,msg: %s ,openMessageID: %s\n", resp.Code, resp.Msg, resp.Data.MessageID)
	}
}

func TestCardSyncButton(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	builder := &message.CardBuilder{}
	//add header
	content := "Please choose color"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")
	builder.AddHRBlock()
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)

	//add block
	builder.AddDIVBlock(nil, []protocol.FieldForm{
		*message.NewField(false, message.NewMDText("**Sync**", nil, nil, nil)),
	}, nil)
	payload1 := make(map[string]string, 0)
	payload1["color"] = "red"
	payload2 := make(map[string]string, 0)
	payload2["color"] = "black"
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewButton(message.NewMDText("red", nil, nil, nil),
			nil, nil, payload1, protocol.DANGER, nil, "syncButton"),
		message.NewButton(message.NewMDText("black", nil, nil, nil),
			nil, nil, payload2, protocol.DANGER, nil, "syncButton"),
	})
	card, err := builder.BuildForm()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, false)
	if err != nil {
		t.Errorf("TestSyncButton failed:%v\n", err)
	} else {
		t.Logf("code:%d ,msg: %s ,openMessageID: %s\n", resp.Code, resp.Msg, resp.Data.MessageID)
	}
}

func TestCardJumpButton(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	builder := &message.CardBuilder{}
	//add header
	content := "Jump demo"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")
	builder.AddHRBlock()
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)
	//add block
	url := "https://www.google.com"
	ext := message.NewJumpButton(message.NewMDText("jump to google", nil, nil, nil), &url, nil, protocol.DEFAULT)
	builder.AddDIVBlock(message.NewMDText("", nil, nil, nil), nil, ext)
	card, err := builder.BuildForm()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, false)
	if err != nil {
		t.Errorf("TestJumpButton failed:%v\n", err)
	} else {
		t.Logf("code:%d ,msg: %s ,openMessageID: %s\n", resp.Code, resp.Msg, resp.Data.MessageID)
	}
}

func TestCardImage(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	builder := &message.CardBuilder{}
	//add header
	content := "cardImage demo"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")
	builder.AddHRBlock()
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)
	//add block
	builder.AddImageBlock(
		message.NewMDText("", nil, nil, nil),
		*message.NewMDText("", nil, nil, nil),
		imageKey)
	builder.AddHRBlock()

	card, err := builder.BuildForm()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, false)
	if err != nil {
		t.Errorf("TestImage failed:%v\n", err)
	} else {
		t.Logf("code:%d ,msg: %s ,openMessageID: %s\n", resp.Code, resp.Msg, resp.Data.MessageID)
	}
}

func TestCardSelectStaticMenu(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	builder := &message.CardBuilder{}
	//add header
	content := "selectStaticMenu demo"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")
	builder.AddHRBlock()
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)
	//add block
	params := map[string]string{}
	content1 := "option1"
	content2 := "option2"
	content3 := "option3"
	option1 := message.NewOption(protocol.TextForm{Tag: "plain_text", Content: &content1}, "option1")
	option2 := message.NewOption(protocol.TextForm{Tag: "plain_text", Content: &content2}, "option2")
	option3 := message.NewOption(protocol.TextForm{Tag: "plain_text", Content: &content3}, "option3")
	options := []protocol.OptionForm{option1, option2, option3}
	initOption := "test"
	menu := message.NewSelectStaticMenu(message.NewMDText("default", nil, nil, nil),
		params, options, &initOption, nil, "selectStaticMenu")
	builder.AddDIVBlock(message.NewMDText("", nil, nil, nil), nil, menu)
	builder.AddHRBlock()
	card, err := builder.BuildForm()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, false)
	if err != nil {
		t.Errorf("TestSelectStaticMenu failed:%v\n", err)
	} else {
		t.Logf("code:%d ,msg: %s ,openMessageID: %s\n", resp.Code, resp.Msg, resp.Data.MessageID)
	}
}

func TestCardSelectPersonMenu(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	builder := &message.CardBuilder{}
	//add header
	content := "selectPersonMenu demo"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")
	builder.AddHRBlock()
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)
	//add block
	params := map[string]string{}
	options := []protocol.OptionForm{{}}
	menu := message.NewSelectPersonMenu(message.NewMDText("default", nil, nil, nil),
		params, options, nil, nil, "selectPersonMenu")
	builder.AddDIVBlock(message.NewMDText("", nil, nil, nil), nil, menu)
	builder.AddHRBlock()
	card, err := builder.BuildForm()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, false)
	if err != nil {
		t.Errorf("TestSelectPersonMenu failed:%v\n", err)
	} else {
		t.Logf("code:%d ,msg: %s ,openMessageID: %s\n", resp.Code, resp.Msg, resp.Data.MessageID)
	}
}

func TestCardOverFlowMenu(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	builder := &message.CardBuilder{}
	//add header
	content := "overFlowMenu demo"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")
	builder.AddHRBlock()
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)
	//add block
	params := map[string]string{}
	content1 := "option1"
	content2 := "option2"
	content3 := "option3"
	option1 := message.NewOption(protocol.TextForm{Tag: "plain_text", Content: &content1}, "option1")
	option2 := message.NewOption(protocol.TextForm{Tag: "plain_text", Content: &content2}, "option2")
	option3 := message.NewOption(protocol.TextForm{Tag: "plain_text", Content: &content3}, "option3")
	options := []protocol.OptionForm{option1, option2, option3}
	menu := message.NewOverflowMenu(params, options, nil, "overFlowMenu")
	builder.AddDIVBlock(message.NewMDText("", nil, nil, nil), nil, menu)
	builder.AddHRBlock()
	card, err := builder.BuildForm()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, false)
	if err != nil {
		t.Errorf("TestOverFlowMenu failed:%v\n", err)
	} else {
		t.Logf("code:%d ,msg: %s ,openMessageID: %s\n", resp.Code, resp.Msg, resp.Data.MessageID)
	}
}

func TestCardDatePicker(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	builder := &message.CardBuilder{}
	//add header
	content := "PickerDate demo"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")
	builder.AddHRBlock()
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)
	//add block
	params := map[string]string{}
	initialDate := "2019-9-1"
	menu := message.NewPickerDate(message.NewMDText("default", nil, nil, nil),
		params, nil, &initialDate, "PickerDate")
	builder.AddDIVBlock(message.NewMDText("", nil, nil, nil), nil, menu)
	builder.AddHRBlock()
	card, err := builder.BuildForm()

	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	resp, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, false)
	if err != nil {
		t.Errorf("TestPickerDate failed:%v\n", err)
	} else {
		t.Logf("code:%d ,msg: %s ,openMessageID: %s\n", resp.Code, resp.Msg, resp.Data.MessageID)
	}
}
