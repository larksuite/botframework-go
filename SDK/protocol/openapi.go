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
	GetChatInfoPath                  OpenApiPath = "/open-apis/chat/v4/info/"
	GetChatListPath                  OpenApiPath = "/open-apis/chat/v4/list/"
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
