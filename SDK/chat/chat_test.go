// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chat_test

import (
	"context"
	"testing"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/chat"
)

func TestGetChatInfo(t *testing.T) {
	c := context.Background()
	InitTestParams()
	appconfig.Init(*appConf)
	resp, err := chat.GetChatInfo(c, tenantKey, appConf.AppID, chatID)
	if err != nil {
		t.Errorf("GetChatInfo failed: %v", err)
	} else {
		t.Logf("GetChatInfo chatInfo[%+v]", resp)
	}
}

func TestGetChatList(t *testing.T) {
	ctx := context.Background()
	InitTestParams()
	appconfig.Init(*appConf)
	resp, err := chat.GetChatList(ctx, tenantKey, appConf.AppID, 100, "")
	if err != nil {
		t.Errorf("GetChatList failed: %v", err)
	} else {
		t.Logf("GetChatList hasMore[%v]pageToken[%s]groupNum[%d]", resp.Data.HasMore, resp.Data.PageToken, len(resp.Data.Groups))
		for k, v := range resp.Data.Groups {
			t.Logf("GetChatList chatList %d:%+v\n", k+1, *v)
		}
	}
}
