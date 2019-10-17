// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

const (
	ISVApp      = "isv"
	InternalApp = "internal"
)

// GetTenantAccessToken Internal App request
type GetTenantAccessTokenInternalReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

// GetTenantAccessToken ISV App request
type GetTenantAccessTokenISVReq struct {
	AppAccessToken string `json:"app_access_token"`
	TenantKey      string `json:"tenant_key"`
}

// GetTenantAccessToken Internal/ISV response
type GetTenantAccessTokenResp struct {
	BaseResponse
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

// GetAppAccessToken Internal request
type GetAppAccessTokenInternalReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

// GetAppAccessToken Internal response
type GetAppAccessTokenInternalResp struct {
	BaseResponse
	Expire            int    `json:"expire"`
	AppAccessToken    string `json:"app_access_token"`
	TenantAccessToken string `json:"tenant_access_token"`
}

// GetAppAccessToken ISV request
type GetAppAccessTokenIsvReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	AppTicket string `json:"app_ticket"`
}

// GetAppAccessToken ISV response
type GetAppAccessTokenIsvResp struct {
	BaseResponse
	Expire         int    `json:"expire"`
	AppAccessToken string `json:"app_access_token"`
}

// AppTicketReq request
type AppTicketReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

// AppTicketResp response
type AppTicketResp struct {
	BaseResponse
}
