// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth_test

import (
	"context"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/chat"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

var (
	once      sync.Once
	appConf   *appconfig.AppConfig
	tenantKey string
	openID    string
)

func initTestParams() {
	once.Do(func() {
		appConf = &appconfig.AppConfig{
			AppID:     os.Getenv("appid"),
			AppSecret: os.Getenv("appsecret"),
			AppType:   os.Getenv("apptype"),
		}

		tenantKey = os.Getenv("tenantkey")
		openID = os.Getenv("openid")

		if appConf.AppType == protocol.ISVApp {
			redisClient := &common.DefaultRedisClient{}
			err := redisClient.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
			if err != nil {
				panic("init DB error")
			}
			auth.InitISVAppTicketManager(auth.NewDefaultAppTicketManager(redisClient))
		}

		appconfig.Init(*appConf)
	})
}

func TestAuthAppToken(t *testing.T) {
	ctx := context.Background()
	initTestParams()

	// from svr
	appToken, err := auth.GetAppAccessToken(ctx, appConf.AppID)
	if err != nil {
		t.Errorf("GetAppAccessToken failed: getToken error[%v]", err)
		return
	}
	t.Logf("GetAppAccessToken Success: appTokenSize[%d]appToken[%s]", len(appToken), appToken)

	// from cache
	appToken2, err := auth.GetAppAccessToken(ctx, appConf.AppID)
	if err != nil {
		t.Errorf("GetAppAccessToken failed: getToken error[%v]", err)
		return
	}
	t.Logf("GetAppAccessToken Success: appTokenSize[%d]appToken2[%s]", len(appToken2), appToken2)

	if appToken != appToken2 {
		t.Errorf("GetAppAccessToken failed: cache data error")
		return
	}

	// set disable
	tokenManager, err := appconfig.GetTokenManager(appConf.AppID)
	if err != nil {
		t.Errorf("GetAppAccessToken failed: getManager error[%v]", err)
		return
	}
	tokenManager.DisableAppAccessToken()

	// from svr
	appToken3, err := auth.GetAppAccessToken(ctx, appConf.AppID)
	if err != nil {
		t.Errorf("GetAppAccessToken failed: getToken error[%v]", err)
		return
	}
	t.Logf("GetAppAccessToken Success: appTokenSize[%d]appToken3[%s]", len(appToken3), appToken3)

	// from cache
	appToken4, err := auth.GetAppAccessToken(ctx, appConf.AppID)
	if err != nil {
		t.Errorf("GetAppAccessToken failed: getToken error[%v]", err)
		return
	}
	t.Logf("GetAppAccessToken Success: appTokenSize[%d]appToken4[%s]", len(appToken4), appToken4)

	if appToken3 != appToken4 {
		t.Errorf("GetAppAccessToken failed: cache data error")
		return
	}

}

func TestAuthTenantToken(t *testing.T) {
	ctx := context.Background()
	initTestParams()

	// from svr
	token, err := auth.GetTenantAccessToken(ctx, tenantKey, appConf.AppID)
	if err != nil {
		t.Errorf("GetTenantAccessToken failed: getToken error[%v]", err)
		return
	}
	t.Logf("GetTenantAccessToken Success: appTokenSize[%d]token[%s]", len(token), token)

	// from cache
	token2, err := auth.GetTenantAccessToken(ctx, tenantKey, appConf.AppID)
	if err != nil {
		t.Errorf("GetTenantAccessToken failed: getToken error[%v]", err)
		return
	}
	t.Logf("GetTenantAccessToken Success: appTokenSize[%d]token2[%s]", len(token2), token2)

	if token != token2 {
		t.Errorf("GetTenantAccessToken failed: cache data error")
		return
	}

	// set disable
	tokenManager, err := appconfig.GetTokenManager(appConf.AppID)
	if err != nil {
		t.Errorf("GetTenantAccessToken failed: getManager error[%v]", err)
		return
	}
	tokenManager.DisableTenantAccessToken(tenantKey)

	// from svr
	token3, err := auth.GetTenantAccessToken(ctx, tenantKey, appConf.AppID)
	if err != nil {
		t.Errorf("GetTenantAccessToken failed: getToken error[%v]", err)
		return
	}
	t.Logf("GetTenantAccessToken Success: appTokenSize[%d]token3[%s]", len(token3), token3)

	// from cache
	token4, err := auth.GetTenantAccessToken(ctx, tenantKey, appConf.AppID)
	if err != nil {
		t.Errorf("GetTenantAccessToken failed: getToken error[%v]", err)
		return
	}
	t.Logf("GetTenantAccessToken Success: appTokenSize[%d]token4[%s]", len(token4), token4)

	if token3 != token4 {
		t.Errorf("GetTenantAccessToken failed: cache data error")
		return
	}

}

func TestAutoDisableTenantToken(t *testing.T) {
	ctx := context.Background()
	initTestParams()

	// check chat
	chatID := ""
	for i := 0; i < 3; i++ {
		if i == 1 {
			invalidTenantToken(appConf.AppID, tenantKey)
		}

		rspChatList, err := chat.GetChatList(ctx, tenantKey, appConf.AppID, 100, "")
		if err != nil {
			t.Errorf("check_autoDisable: GetChatList failed, error[%v]", err)
			continue
		}

		for _, v := range rspChatList.Data.Groups {
			chatID = v.ChatID
			break
		}
	}

	for i := 0; i < 3; i++ {
		if i == 1 {
			invalidTenantToken(appConf.AppID, tenantKey)
		}

		_, err := chat.GetChatInfo(ctx, tenantKey, appConf.AppID, chatID)
		if err != nil {
			t.Errorf("check_autoDisable: GetChatInfo failed, error[%v]", err)
			continue
		}
	}

	// check upload image
	for i := 0; i < 3; i++ {
		if i == 1 {
			invalidTenantToken(appConf.AppID, tenantKey)
		}

		path := "../../demo/source/lark" + strconv.Itoa(i) + ".jpg"
		_, err := message.GetImageKey(ctx, tenantKey, appConf.AppID, "", path)
		if err != nil {
			t.Errorf("check_autoDisable: GetImageKey failed, error[%v]", err)
			continue
		}
	}

	// check send message
	user := &protocol.UserInfo{
		ID:   openID,
		Type: protocol.UserTypeOpenID,
	}

	for i := 0; i < 3; i++ {
		if i == 1 {
			invalidTenantToken(appConf.AppID, tenantKey)
		}

		text := "sendTextMessage test(~_~)"
		_, err := message.SendTextMessage(ctx, tenantKey, appConf.AppID, user, "", text)
		if err != nil {
			t.Errorf("check_autoDisable: SendTextMessage failed, error[%v]", err)
			continue
		}
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

	//generate card
	card, _ := builder.BuildForm()

	for i := 0; i < 3; i++ {
		if i == 1 {
			invalidTenantToken(appConf.AppID, tenantKey)
		}

		_, err := message.SendCardMessage(ctx, tenantKey, appConf.AppID, user, "", *card, true)
		if err != nil {
			t.Errorf("check_autoDisable: SendCardMessage failed, error[%v]", err)
			continue
		}
	}

	// check batch send message
	info := &protocol.BatchBaseInfo{
		DepartmentIDs: []string{},
		OpenIDs:       []string{openID},
		UserIDs:       []string{},
	}

	for i := 0; i < 3; i++ {
		if i == 1 {
			invalidTenantToken(appConf.AppID, tenantKey)
		}

		text := "sendTextMessage test(~_~)"
		_, err := message.SendTextMessageBatch(ctx, tenantKey, appConf.AppID, info, "", text)
		if err != nil {
			t.Errorf("check_autoDisable: SendTextMessageBatch failed, error[%v]", err)
			continue
		}
	}

	for i := 0; i < 3; i++ {
		if i == 1 {
			invalidTenantToken(appConf.AppID, tenantKey)
		}

		_, err := message.SendCardMessageBatch(ctx, tenantKey, appConf.AppID, info, "", *card, true)
		if err != nil {
			t.Errorf("check_autoDisable: SendCardMessageBatch failed, error[%v]", err)
			continue
		}
	}

}

func invalidTenantToken(appID, tenantKey string) error {
	token, err := auth.GetTenantAccessToken(context.TODO(), tenantKey, appID)
	if err != nil {
		return err
	}

	tokenManager, err := appconfig.GetTokenManager(appConf.AppID)
	if err != nil {
		return err
	}

	invalidToken := []byte(token)[0 : len(token)-4]
	tokenManager.SetTenantAccessToken(tenantKey, string(invalidToken), 3600)

	return nil
}

func invalidAppToken(appID string) error {
	token, err := auth.GetAppAccessToken(context.TODO(), appID)
	if err != nil {
		return err
	}

	tokenManager, err := appconfig.GetTokenManager(appConf.AppID)
	if err != nil {
		return err
	}

	invalidToken := []byte(token)[0 : len(token)-4]
	tokenManager.SetAppAccessToken(string(invalidToken), 3600)

	return nil
}
