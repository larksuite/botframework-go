// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

type MiniProgramLoginByAppTokenRequest struct {
	Code string `json:"code"`
}

type MiniProgramLoginByAppTokenResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`

	Data struct {
		MiniProgramUserToken
	} `json:"data"`
}

func GenMiniProgramLoginByIDSecretRequest(code, appID, appSecret string) map[string]string {
	return map[string]string{
		"code":   code,
		"appid":  appID,
		"secret": appSecret,
	}
}

type MiniProgramLoginByIDSecretResponse struct {
	Code int    `json:"error"`
	Msg  string `json:"message"`

	MiniProgramUserToken
}

type MiniProgramUserToken struct {
	AccessToken  string `json:"access_token"` // user_access_token
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TenantKey    string `json:"tenant_key"`
	OpenID       string `json:"open_id"`
	UnionID      string `json:"union_id"`
	SessionKey   string `json:"session_key"`
}

const (
	GrantTypeAuthCode     = "authorization_code"
	GrantTypeRefreshToken = "refresh_token"
)

type OpenSSOTokenRequest struct {
	//(app_id + app_secret)/app_access_token, two ways to choose one. app_access_token is recommended
	AppID          string `json:"app_id"`
	AppSecret      string `json:"app_secret"`
	AppAccessToken string `json:"app_access_token"`
	GrantType      string `json:"grant_type"`    // authorization_code / refresh_token
	Code           string `json:"code"`          // required when GrantType=authorization_code
	RefreshToken   string `json:"refresh_token"` // required when GrantType=refresh_token
}

type OpenSSOTokenResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`

	UserTokenInfo
}

type UserTokenInfo struct {
	AccessToken  string `json:"access_token"` // user_access_token
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TenantKey    string `json:"tenant_key"`
	OpenID       string `json:"open_id"`
	EmployeeID   string `json:"employee_id"`
	Name         string `json:"name"`
	ENName       string `json:"en_name"`
	AvatarURL    string `json:"avatar_url"`
}
