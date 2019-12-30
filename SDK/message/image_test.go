// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package message_test

import (
	"context"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/message"
)

var (
	once      sync.Once
	appConf   *appconfig.AppConfig
	tenantKey string
	chatID    string
	openID    string
	userID    string
	imageKey  string
	messageID string
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

		tenantKey = os.Getenv("tenantkey")
		chatID = os.Getenv("chatid")
		openID = os.Getenv("openid")
		userID = os.Getenv("userid")

		imageKey = os.Getenv("imagekey")

		messageID = os.Getenv("messageid")

		appconfig.Init(*appConf)
	})
}

func TestGetImageKey(t *testing.T) {
	c := context.Background()
	InitTestParams()

	// by path
	path := "../../demo/source/lark0.jpg"
	key, err := message.GetImageKey(c, tenantKey, appConf.AppID, "", path)
	if err != nil {
		t.Errorf("GetImageKeyByPath failed: %v", err)
	} else {
		t.Logf("GetImageKeyByPath: image_key = %s", key)
	}

	// by url
	url := "https://s0.pstatp.com/ee/lark-open/web/static/apply.226f11cb.png"
	key, err = message.GetImageKey(c, tenantKey, appConf.AppID, url, "")
	if err != nil {
		t.Errorf("GetImageKeyByURL failed: %v", err)
	} else {
		t.Logf("GetImageKeyByURL: image_key = %s", key)
	}
}

func TestGetImageBinData(t *testing.T) {
	c := context.Background()
	InitTestParams()

	data, err := message.GetImageBinData(c, tenantKey, appConf.AppID, imageKey)
	if err != nil {
		t.Errorf("GetImageBinData failed: %v", err)
	} else {
		t.Logf("GetImageBinData: succ")
	}

	err = ioutil.WriteFile("./temp.jpg", data, 0644)
	if err != nil {
		t.Errorf("WriteFile failed: %v", err)
	} else {
		t.Logf("WriteFile: succ")
	}

}
