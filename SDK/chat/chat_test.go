// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chat_test

import (
	"context"
	"testing"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/chat"
	"github.com/larksuite/botframework-go/SDK/protocol"
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

func TestUpdateChatInfo(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	appconfig.Init(*appConf)

	newName := "Update Name"
	request := &protocol.UpdateChatInfoRequest{
		ChatID:        chatID,
		OwnerUserID:   nil,
		OwnerOpenID:   nil,
		Name:          &newName,
		ChatI18nNames: nil,
	}

	resp, err := chat.UpdateChatInfo(ctx, tenantKey, appConf.AppID, request)
	if err != nil {
		t.Errorf("UpdateChatInfo failed: %v", err)
	} else {
		t.Logf("UpdateChatInfo Success, chatid: [%+v]", resp)
	}
}

func TestCreateChat(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	appconfig.Init(*appConf)

	userids := []string{}
	openids := []string{openID}
	request := &protocol.CreateChatRequest{
		Name:          "test group",
		Description:   "test group description",
		UserIDs:       userids,
		OpenIDs:       openids,
		ChatI18nNames: nil,
	}

	resp, err := chat.CreateChat(ctx, tenantKey, appConf.AppID, request)
	if err != nil {
		t.Errorf("CreateChat failed: %v", err)
	} else {
		t.Logf("CreateChat Success, chatid: [%+v]", resp)
	}
}

func TestAddUserToChat(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	appconfig.Init(*appConf)

	userids := []string{}
	openids := []string{openID}
	request := &protocol.AddUserToChatRequest{
		ChatID:  chatID,
		UserIDs: userids,
		OpenIDs: openids,
	}

	resp, err := chat.AddUserToChat(ctx, tenantKey, appConf.AppID, request)
	if err != nil {
		t.Errorf("AddUserToChat failed: %v", err)
	} else {
		t.Logf("AddUserToChat Success [%+v]", resp)
	}
}

func TestDeleteUserFromChat(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	appconfig.Init(*appConf)

	userids := []string{}
	openids := []string{openID}
	request := &protocol.DeleteUserFromChatRequest{
		ChatID:  chatID,
		UserIDs: userids,
		OpenIDs: openids,
	}

	resp, err := chat.DeleteUserFromChat(ctx, tenantKey, appConf.AppID, request)
	if err != nil {
		t.Errorf("DeleteUserFromChat failed: %v", err)
	} else {
		t.Logf("DeleteUserFromChat Success [%+v]", resp)
	}
}

func TestDisbandChat(t *testing.T) {
	ctx := context.Background()
	InitTestParams()

	appconfig.Init(*appConf)

	request := &protocol.DisbandChatRequest{
		ChatID: chatID,
	}

	resp, err := chat.DisbandChat(ctx, tenantKey, appConf.AppID, request)
	if err != nil {
		t.Errorf("DisbandChat failed: %v", err)
	} else {
		t.Logf("DisbandChat Success [%+v]", resp)
	}
}
