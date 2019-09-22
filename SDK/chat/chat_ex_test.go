// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chat_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/chat"
)

var (
	once      sync.Once
	appConf   *appconfig.AppConfig
	tenantKey string
	chatID    string
	openID    string
	userID    string
)

func InitTestParams() {
	once.Do(func() {
		appConf = &appconfig.AppConfig{
			AppID:       os.Getenv("appid"),
			AppSecret:   os.Getenv("appsecret"),
			VerifyToken: os.Getenv("verifytoken"),
			EncryptKey:  os.Getenv("encryptkey"),
			AppType:     os.Getenv("apptype"),
		}
		// can get from other way like redis
		tenantKey = os.Getenv("tenantkey")
		chatID = os.Getenv("chatid")
		openID = os.Getenv("openid")
		userID = os.Getenv("userid")

		appconfig.Init(*appConf)
	})
}

func TestCheckUserInGroup(t *testing.T) {
	c := context.Background()
	InitTestParams()

	inGroup, err := chat.CheckOpenIDInGroup(c, tenantKey, appConf.AppID, chatID, openID)
	if err != nil {
		t.Errorf("CheckOpenIDInGroup failed: %v", err)
	} else {
		t.Logf("CheckOpenIDInGroup: %v", inGroup)
	}

	inGroup, err = chat.CheckUserIDInGroup(c, tenantKey, appConf.AppID, chatID, userID)
	if err != nil {
		t.Errorf("CheckUserIDInGroup failed: %v", err)
	} else {
		t.Logf("CheckUserIDInGroup: %v", inGroup)
	}
}

func TestCheckBotInGroup(t *testing.T) {
	c := context.Background()
	InitTestParams()

	inGroup, err := chat.CheckBotInGroup(c, tenantKey, appConf.AppID, chatID)
	if err != nil {
		t.Errorf("CheckBotInGroup failed: %v", err)
	} else {
		t.Logf("CheckBotInGroup: %v", inGroup)
	}
}

func TestCheckUserBotInSameGroup(t *testing.T) {
	c := context.Background()
	InitTestParams()

	inGroup, err := chat.CheckOpenIDBotInSameGroup(c, tenantKey, appConf.AppID, chatID, openID)
	if err != nil {
		t.Errorf("CheckOpenIDBotInSameGroup failed: %v", err)
	} else {
		t.Logf("CheckOpenIDBotInSameGroup: %v", inGroup)
	}

	inGroup, err = chat.CheckUserIDBotInSameGroup(c, tenantKey, appConf.AppID, chatID, userID)
	if err != nil {
		t.Errorf("CheckUserIDBotInSameGroup failed: %v", err)
	} else {
		t.Logf("CheckUserIDBotInSameGroup: %v", inGroup)
	}
}
