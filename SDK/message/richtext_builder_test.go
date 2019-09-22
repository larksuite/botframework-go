// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package message_test

import (
	"encoding/json"
	"testing"

	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

func TestBuildRichtextReq(t *testing.T) {
	postForm := make(map[protocol.Language]*protocol.RichTextForm)

	// en-us
	titleUS := "this is a title"
	contentUS := message.NewRichTextContent()

	// first line
	contentUS.AddElementBlock(
		message.NewTextTag("first line :", true, 1),
		message.NewATag("hyperlinks", true, "https://www.feishu.cn"),
		message.NewAtTag("username", "userid"),
	)

	// second line
	contentUS.AddElementBlock(
		message.NewTextTag("second line :", true, 1),
		message.NewTextTag("text test", true, 1),
	)

	postForm[protocol.EnUS] = message.NewRichTextForm(&titleUS, contentUS)

	// zh-cn
	titleCN := "这是一个标题"
	contentCN := message.NewRichTextContent()

	// first line
	contentCN.AddElementBlock(
		message.NewTextTag("第一行 :", true, 1),
		message.NewATag("超链接", true, "https://www.feishu.cn"),
		message.NewAtTag("username", "userid"),
	)

	// second line
	contentCN.AddElementBlock(
		message.NewTextTag("第二行 :", true, 1),
		message.NewTextTag("文本测试", true, 1),
	)

	postForm[protocol.ZhCN] = message.NewRichTextForm(&titleCN, contentCN)

	post := map[string]*protocol.RichTextForm{}
	for k, v := range postForm {
		post[k.String()] = v
	}

	bytes, err := json.Marshal(post)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(bytes))
}
