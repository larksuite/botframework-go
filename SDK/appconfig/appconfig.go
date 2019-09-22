// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package appconfig

import "fmt"

var (
	appConfMap  = make(map[string]AppConfig)
	appTokenMap = make(map[string]*AppToken)
)

type AppConfig struct {
	AppID       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	VerifyToken string `json:"verify_token"`
	EncryptKey  string `json:"encrypt_key"`
	AppType     string `json:"app_type"`
}

type AppToken struct {
	AppAccessToken    *AppAccessTokenCache
	TenantAccessToken map[string]*TenantAccessTokenCache //tenantKey->tenantAccessToken
}

type AppAccessTokenCache struct {
	Token  string
	Expire int64
}

type TenantAccessTokenCache struct {
	Token  string
	Expire int64
}

func Init(appConfs ...AppConfig) {
	for _, v := range appConfs {
		appConfMap[v.AppID] = v

		appTokenMap[v.AppID] = &AppToken{
			AppAccessToken:    &AppAccessTokenCache{},
			TenantAccessToken: make(map[string]*TenantAccessTokenCache),
		}
	}
}

func GetConfig(appID string) (AppConfig, error) {
	appConf, ok := appConfMap[appID]
	if !ok {
		return AppConfig{}, fmt.Errorf("getAppConfig: cannot find appConfig, appid[%s]confSize[%d]", appID, len(appConfMap))
	}

	return appConf, nil
}

func GetToken(appID string) (*AppToken, error) {
	appToken, ok := appTokenMap[appID]
	if !ok {
		return nil, fmt.Errorf("getAppToken: cannot find appToken, appid[%s]tokenSize[%d]", appID, len(appTokenMap))
	}

	return appToken, nil
}
