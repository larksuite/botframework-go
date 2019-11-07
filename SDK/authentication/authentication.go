// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authentication

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Authentication interface {
	Init(manager SessionManager, validPeriod time.Duration)

	Login(ctx context.Context, code string, appID string, requestHost string) (map[string]*http.Cookie, error)
	Auth(ctx context.Context, sessionKey string) error
	Logout(ctx context.Context, appID string, requestHost string) (map[string]*http.Cookie, error)

	// auth cookie level, DomainLevelZero is default
	SetCookieDomainLevel(cookieLevel CookieDomainLevel)
	GetCookieDomainLevel() CookieDomainLevel

	// get valid period
	GetValidPeriod() time.Duration

	// get session manager
	GetSessionManager() SessionManager
}

type CookieDomainLevel int

const (
	DomainLevelZero CookieDomainLevel = 0 //current host
	DomainLevelOne  CookieDomainLevel = 1 //First-level domain name
	DomainLevelTwo  CookieDomainLevel = 2 //Second-level domain name
)

func GetAuthCookieDomain(requestHost string, cookieLevel CookieDomainLevel) string {
	if requestHost == "" {
		return ""
	}

	names := strings.Split(requestHost, ".")
	namesLen := len(names)

	switch cookieLevel {
	case DomainLevelOne:
		if namesLen <= 2 {
			return requestHost
		} else {
			return fmt.Sprintf("%s.%s", names[namesLen-2], names[namesLen-1])
		}
	case DomainLevelTwo:
		if namesLen <= 3 {
			return requestHost
		} else {
			return fmt.Sprintf("%s.%s.%s", names[namesLen-3], names[namesLen-2], names[namesLen-1])
		}
	default:
		return ""
	}
}

type SessionManager interface {
	SetEncryptKey(encryptKey string)
	GetEncryptKey() string

	GenerateSessionKeyName(appID string) string
	GenerateSessionKey() string

	SetAuthUserInfo(authUser *AuthUserInfo, validPeriod time.Duration) (string, error) // return sessionKey, err
	GetAuthUserInfo(sessionKey string) (*AuthUserInfo, error)
}

type AuthUserInfo struct {
	Token TokenInfo
	User  UserInfo
	Extra map[string]string
}

type TokenInfo struct {
	AccessToken  string `json:"access_token"` // user_access_token
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	TenantKey  string `json:"tenant_key"`
	OpenID     string `json:"open_id"`
	EmployeeID string `json:"employee_id"`
}
