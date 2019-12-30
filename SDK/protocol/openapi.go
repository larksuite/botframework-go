// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

type OpenApiPath string

const (
	GetAppAccessTokenInternalPath    OpenApiPath = "/open-apis/auth/v3/app_access_token/internal/"
	GetAppAccessTokenIsvPath         OpenApiPath = "/open-apis/auth/v3/app_access_token/"
	GetTenantAccessTokenInternalPath OpenApiPath = "/open-apis/auth/v3/tenant_access_token/internal/"
	GetTenantAccessTokenIsvPath      OpenApiPath = "/open-apis/auth/v3/tenant_access_token/"
	ResendAppTicketPath              OpenApiPath = "/open-apis/auth/v3/app_ticket/resend"
	SendMessagePath                  OpenApiPath = "/open-apis/message/v4/send/"
	SendMessageBatchPath             OpenApiPath = "/open-apis/message/v4/batch_send/"
	UploadImagePath                  OpenApiPath = "/open-apis/image/v4/upload/"
	GetImagePath                     OpenApiPath = "/open-apis/image/v4/get"
	GetChatInfoPath                  OpenApiPath = "/open-apis/chat/v4/info/"
	GetChatListPath                  OpenApiPath = "/open-apis/chat/v4/list/"
	UpdateChatInfoPath               OpenApiPath = "/open-apis/chat/v4/update/"
	CreateChatPath                   OpenApiPath = "/open-apis/chat/v4/create/"
	AddUserToChatPath                OpenApiPath = "/open-apis/chat/v4/chatter/add/"
	DeleteUserFromChatPath           OpenApiPath = "/open-apis/chat/v4/chatter/delete/"
	DisbandChatPath                  OpenApiPath = "/open-apis/chat/v4/disband/"
	CardUpdatePath                   OpenApiPath = "/open-apis/interactive/v1/card/update"
	MPValidateByAppTokenPath         OpenApiPath = "/open-apis/mina/v2/tokenLoginValidate"   //mini programe login validate, ExchangeToken
	MPValidateByIDSecretPath         OpenApiPath = "/open-apis/mina/loginValidate"           //mini programe login validate, ExchangeToken
	OpenSSOValidatePath              OpenApiPath = "/connect/qrconnect/oauth2/access_token/" //open sso login validate, ExchangeToken/RefreshToken
	OpenSSOGetCodePath               OpenApiPath = "/connect/qrconnect/page/sso/"            //open sso, GetCode
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const (
	ErrAppTicketNil              = 10003
	ErrAppTicketInvalid          = 10012
	ErrTenantAccessTokenInvalid  = 99991663
	ErrAppAccessTokenInvalid     = 99991664
	ErrMinaAppAccessTokenInvalid = 10202
)
