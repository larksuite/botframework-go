// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authentication

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/larksuite/botframework-go/SDK/common"
	uuid "github.com/satori/go.uuid"
)

type defaultSessionManager struct {
	EncryptKey string
	Client     common.DBClient
}

// NewDefaultSessionManager demo:
// client := &common.DefaultRedisClient{}
// err := client.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
// if err != nil {
// 	return fmt.Errorf("init redis error[%v]", err)
// }
// manager := authentication.NewDefaultSessionManager("DojK2hs*790(", client)
func NewDefaultSessionManager(encryptKey string, client common.DBClient) *defaultSessionManager {
	manager := &defaultSessionManager{
		EncryptKey: encryptKey,
		Client:     client,
	}

	return manager
}

func (d *defaultSessionManager) SetEncryptKey(encryptKey string) {
	d.EncryptKey = encryptKey
}

func (d *defaultSessionManager) GetEncryptKey() string {
	return d.EncryptKey
}

func (d *defaultSessionManager) GenerateSessionKeyName(appID string) string {
	prefix := "cli_"
	if strings.HasPrefix(appID, prefix) {
		appID = appID[len(prefix):]
	}

	return fmt.Sprintf("bframewk-session-%s", appID)
}

func (d *defaultSessionManager) GenerateSessionKey() string {
	newUUID := uuid.NewV4()
	return fmt.Sprintf("bframewk-w04j5mfw-%s-%s", newUUID.String(), strconv.FormatInt(time.Now().Unix(), 36))
}

func (d *defaultSessionManager) SetAuthUserInfo(authUser *AuthUserInfo, validPeriod time.Duration) (string, error) {
	sessionKey := d.GenerateSessionKey()

	sessionValue, err := common.EncryptAes(authUser, sessionKey+d.GetEncryptKey())
	if err != nil {
		return "", fmt.Errorf("auth encrypt err[%v]", err)
	}

	err = d.Client.Set(sessionKey, sessionValue, validPeriod)
	if err != nil {
		return "", fmt.Errorf("set auth err[%v]", err)
	}

	return sessionKey, nil
}

func (d *defaultSessionManager) GetAuthUserInfo(sessionKey string) (*AuthUserInfo, error) {
	sessionValue, err := d.Client.Get(sessionKey)
	if err != nil {
		return nil, fmt.Errorf("get auth err[%v]", err)
	}

	authInfo := &AuthUserInfo{}

	err = common.DecryptAes(sessionValue, sessionKey+d.GetEncryptKey(), authInfo)
	if err != nil {
		return nil, fmt.Errorf("auth decrypt err[%v]", err)
	}

	return authInfo, nil
}
