// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package appconfig

import (
	"fmt"
	"sync"
	"time"
)

var (
	appConfMap = make(map[string]AppConfig)

	appTokenMap = make(map[string]*AppTokenManager)
)

const (
	ExpireInterval = 300 // 300 seconds
)

type AppConfig struct {
	AppID       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	VerifyToken string `json:"verify_token"`
	EncryptKey  string `json:"encrypt_key"`
	AppType     string `json:"app_type"`
}

type AppTokenManager struct {
	AppAccessToken    *AppAccessTokenCache
	TenantAccessToken map[string]*TenantAccessTokenCache //tenantKey->tenantAccessToken
	rwMuApp           sync.RWMutex
	rwMuTenant        sync.RWMutex
}

type AppAccessTokenCache struct {
	Token  string
	Expire int64
}

type TenantAccessTokenCache struct {
	Token  string
	Expire int64
}

func (a *AppTokenManager) GetAppAccessToken() (string, error) {
	a.rwMuApp.RLock()
	defer a.rwMuApp.RUnlock()

	if a.AppAccessToken != nil && a.AppAccessToken.Token != "" && a.AppAccessToken.Expire > time.Now().Unix() {
		return a.AppAccessToken.Token, nil
	}

	return "", fmt.Errorf("cannot find app access token")
}

func (a *AppTokenManager) SetAppAccessToken(appAccessToken string, expireSecond int) error {
	a.rwMuApp.Lock()
	defer a.rwMuApp.Unlock()

	if a.AppAccessToken == nil {
		a.AppAccessToken = new(AppAccessTokenCache)
	}
	a.AppAccessToken.Token = appAccessToken
	a.AppAccessToken.Expire = time.Now().Unix() + int64(expireSecond-ExpireInterval)

	return nil
}

func (a *AppTokenManager) DisableAppAccessToken() error {
	a.rwMuApp.Lock()
	defer a.rwMuApp.Unlock()

	if a.AppAccessToken != nil {
		a.AppAccessToken.Expire = 0
	}

	return nil
}

func (a *AppTokenManager) GetTenantAccessToken(tenantKey string) (string, error) {
	a.rwMuTenant.RLock()
	defer a.rwMuTenant.RUnlock()

	tcToken, ok := a.TenantAccessToken[tenantKey]
	if ok && tcToken != nil && tcToken.Token != "" && tcToken.Expire > time.Now().Unix() {
		return tcToken.Token, nil
	}

	return "", fmt.Errorf("cannot find tenant access token")
}

func (a *AppTokenManager) SetTenantAccessToken(tenantKey string, tenantAccessToken string, expireSecond int) error {
	a.rwMuTenant.Lock()
	defer a.rwMuTenant.Unlock()

	if a.TenantAccessToken == nil {
		a.TenantAccessToken = make(map[string]*TenantAccessTokenCache)
	}
	if a.TenantAccessToken[tenantKey] == nil {
		a.TenantAccessToken[tenantKey] = new(TenantAccessTokenCache)
	}

	a.TenantAccessToken[tenantKey].Token = tenantAccessToken
	a.TenantAccessToken[tenantKey].Expire = time.Now().Unix() + int64(expireSecond-ExpireInterval)

	return nil
}

func (a *AppTokenManager) DisableTenantAccessToken(tenantKey string) error {
	a.rwMuTenant.Lock()
	defer a.rwMuTenant.Unlock()

	if a.TenantAccessToken != nil && a.TenantAccessToken[tenantKey] != nil {
		a.TenantAccessToken[tenantKey].Expire = 0
	}

	return nil
}

func Init(appConfs ...AppConfig) {
	for _, v := range appConfs {
		appConfMap[v.AppID] = v

		appTokenMap[v.AppID] = &AppTokenManager{
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

func GetTokenManager(appID string) (*AppTokenManager, error) {
	tokenManager, ok := appTokenMap[appID]
	if !ok {
		return nil, fmt.Errorf("getAppToken: cannot find tokenManager, appid[%s]tokenSize[%d]", appID, len(appTokenMap))
	}

	return tokenManager, nil
}
