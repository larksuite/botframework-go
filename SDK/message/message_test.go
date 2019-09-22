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
	url := "http://a.hiphotos.baidu.com/image/pic/item/838ba61ea8d3fd1fc9c7b6853a4e251f94ca5f46.jpg"
	resp, err := message.SendImageMessage(c, tenantKey, appConf.AppID, user, "", url, "", "")
	if err != nil {
		t.Errorf("SendImageMessage: failed err[%v]", err)
	} else {
		t.Logf("SendImageMessage: succ messageID[%s]", resp.Data.MessageID)
	}
	//by path
	path := "/tmp/test.png"
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

	text := "this is a test"
	elem := &protocol.RichTextElementForm{
		Tag:  "text",
		Text: &text,
	}
	content := &protocol.RichTextContent{}
	content.AddElementBlock(elem)
	postForm := map[protocol.Language]*protocol.RichTextForm{
		protocol.EnUS: &protocol.RichTextForm{
			Title:   "THIS IS A TITLE",
			Content: content,
		},
	}

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
	resp, err := message.SendCardMessage(c, tenantKey, appConf.AppID, user, "", card, true)
	if err != nil {
		t.Errorf("SendCardMessage: failed err[%v]", err)
	} else {
		t.Logf("SendCardMessage: succ messageID[%s]", resp.Data.MessageID)
	}
}

//need to get access
func TestSendTextMessageBatch(t *testing.T) {
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
func TestSendImageMessageBatch(t *testing.T) {
	c := context.Background()
	InitTestParams()

	info := &protocol.BatchBaseInfo{
		DepartmentIDs: []string{},
		OpenIDs:       []string{openID},
		UserIDs:       []string{},
	}

	//choose url,path or imageKey
	//by url
	url := "http://a.hiphotos.baidu.com/image/pic/item/838ba61ea8d3fd1fc9c7b6853a4e251f94ca5f46.jpg"
	resp, err := message.SendImageMessageBatch(c, tenantKey, appConf.AppID, info, "", url, "", "")
	if err != nil {
		t.Errorf("SendImageMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendImageMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}

	//by path
	path := "/tmp/test.png"
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
func TestSendRichTextMessageBatch(t *testing.T) {
	c := context.Background()
	InitTestParams()
	info := &protocol.BatchBaseInfo{
		DepartmentIDs: []string{},
		OpenIDs:       []string{openID},
		UserIDs:       []string{},
	}
	text := "this is a test"
	elem := &protocol.RichTextElementForm{
		Tag:      "text",
		Text:     &text,
		ImageKey: imageKey,
	}
	content := &protocol.RichTextContent{}
	content.AddElementBlock(elem)
	postForm := map[protocol.Language]*protocol.RichTextForm{
		protocol.EnUS: &protocol.RichTextForm{
			Title:   "THIS IS A TITLE",
			Content: content,
		},
	}
	resp, err := message.SendRichTextMessageBatch(c, tenantKey, appConf.AppID, info, "", postForm)
	if err != nil {
		t.Errorf("SendRichTextMessageBatch: failed err[%v]", err)
	} else {
		t.Logf("SendRichTextMessageBatch: succ messageID[%s]", resp.Data.MessageID)
	}
}

//need to get access
func TestSendShareChatMessageBatch(t *testing.T) {
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
func TestSendCardMessageBatch(t *testing.T) {
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
