// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

var (
	appTicketManager TicketManager
)

// TicketManager set/get your app-ticket。
type TicketManager interface {
	SetAppTicket(ctx context.Context, appID, appTicket string) error
	GetAppTicket(ctx context.Context, appID string) (string, error)
}

func InitISVAppTicketManager(ticketManager TicketManager) error {
	if ticketManager == nil {
		return errors.New("param ticketManager is nil")
	}

	appTicketManager = ticketManager
	return nil
}

func GetTenantAccessToken(ctx context.Context, tenantKey, appID string) (string, error) {
	tokenManager, err := appconfig.GetTokenManager(appID)
	if err != nil {
		return "", common.ErrAppTokenNotFound.ErrorWithExtErr(err)
	}

	cacheTenantAccessToken, err := tokenManager.GetTenantAccessToken(tenantKey)
	if err == nil && cacheTenantAccessToken != "" {
		return cacheTenantAccessToken, nil
	}

	appInfo, err := appconfig.GetConfig(appID)
	if err != nil {
		return "", common.ErrAppConfNotFound.ErrorWithExtErr(err)
	}

	var respBodyObj *protocol.GetTenantAccessTokenResp

	if appInfo.AppType != protocol.ISVApp {
		respBodyObj, err = getInternalTenantAccessToken(ctx, appInfo.AppID, appInfo.AppSecret)
		if err != nil {
			return "", common.ErrGetInternalTenantAccessToken.ErrorWithExtErr(err)
		}
	} else {
		appAccessToken, err := GetAppAccessToken(ctx, appID)
		if err != nil {
			return "", common.ErrGetAppAccessToken.ErrorWithExtErr(err)
		}
		respBodyObj, err = getIsvTenantAccessToken(ctx, tenantKey, appAccessToken)
		if err != nil {
			return "", common.ErrGetISVTenantAccessToken.ErrorWithExtErr(err)
		}
	}

	if respBodyObj == nil {
		return "", common.ErrRespDataIsNil.Error()
	}

	tokenManager.SetTenantAccessToken(tenantKey, respBodyObj.TenantAccessToken, respBodyObj.Expire)

	common.Logger(ctx).Infof("SDK-Get-TenantAccessToken: getFormSvr success appID[%s]tenantKey[%s]tokenSize[%d]expireSecond[%d]",
		appID, tenantKey, len(respBodyObj.TenantAccessToken), respBodyObj.Expire)

	return respBodyObj.TenantAccessToken, nil
}

// ReSendAppTicket app-ticket will be pushed to this service when call this function
func ReSendAppTicket(ctx context.Context, appID, appSecret string) error {
	reqData := &protocol.AppTicketReq{
		AppID:     appID,
		AppSecret: appSecret,
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.ResendAppTicketPath, nil, reqData)
	if err != nil {
		return common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.AppTicketResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		return common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return nil
}

// RefreshAppTicket set app ticket
func RefreshAppTicket(ctx context.Context, data []byte) error {
	if appTicketManager == nil {
		return common.ErrTicketManagerNotInit.Error()
	}

	var appTicket protocol.AppTicketEvent
	err := json.Unmarshal(data, &appTicket)
	if err != nil {
		return common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	err = appTicketManager.SetAppTicket(ctx, appTicket.AppID, appTicket.AppTicket)
	if err != nil {
		return common.ErrSetAppTicketFailed.ErrorWithExtErr(err)
	}

	return nil
}

func GetAppAccessToken(ctx context.Context, appID string) (string, error) {
	tokenManager, err := appconfig.GetTokenManager(appID)
	if err != nil {
		return "", common.ErrAppTokenNotFound.ErrorWithExtErr(err)
	}

	cacheAppAccessToken, err := tokenManager.GetAppAccessToken()
	if err == nil && cacheAppAccessToken != "" {
		return cacheAppAccessToken, nil
	}

	appInfo, err := appconfig.GetConfig(appID)
	if err != nil {
		return "", common.ErrAppConfNotFound.ErrorWithExtErr(err)
	}

	var appAccessToken string
	var expireSecond int
	if appInfo.AppType != protocol.ISVApp {
		rspData, err := getInternalAppAccessToken(ctx, appInfo.AppID, appInfo.AppSecret)
		if err != nil {
			return "", common.ErrGetInternalAppAccessToken.ErrorWithExtErr(err)
		}

		appAccessToken = rspData.AppAccessToken
		expireSecond = rspData.Expire
	} else {
		appTicket, err := appTicketManager.GetAppTicket(ctx, appID)
		if err != nil || appTicket == "" {
			var resultReSend string
			errReSend := ReSendAppTicket(ctx, appInfo.AppID, appInfo.AppSecret) // donot find appTicket
			if errReSend != nil {
				resultReSend = fmt.Sprintf("ReSendAppTicketError[%v]", errReSend)
			} else {
				resultReSend = "ReSendAppTicketSucc"
			}

			return "", common.ErrAppTicketNotFound.ErrorWithExtStr(fmt.Sprintf("getTicketError[%v] or ticketIsEmpty, %s", err, resultReSend))
		}

		rspData, err := getIsvAppAccessToken(ctx, appInfo.AppID, appInfo.AppSecret, appTicket)
		if err != nil {
			return "", common.ErrGetISVAppAccessToken.ErrorWithExtErr(err)
		}

		appAccessToken = rspData.AppAccessToken
		expireSecond = rspData.Expire
	}

	tokenManager.SetAppAccessToken(appAccessToken, expireSecond)

	common.Logger(ctx).Infof("SDK-Get-AppAccessToken: getFormSvr success appID[%s]tokenSize[%d]expireSecond[%d]",
		appID, len(appAccessToken), expireSecond)

	return appAccessToken, nil
}

func getInternalTenantAccessToken(ctx context.Context, appID, appSecret string) (*protocol.GetTenantAccessTokenResp, error) {
	reqData := &protocol.GetTenantAccessTokenInternalReq{
		AppID:     appID,
		AppSecret: appSecret,
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.GetTenantAccessTokenInternalPath, nil, reqData)
	if err != nil {
		return nil, fmt.Errorf("doHttpOApiError[%v]", err)
	}

	rspData := &protocol.GetTenantAccessTokenResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes))
	}

	if rspData.Code != 0 {
		return nil, fmt.Errorf("oapiReturnError[code:%d msg:%s]", rspData.Code, rspData.Msg)
	}
	return rspData, nil
}

func getIsvTenantAccessToken(ctx context.Context, tenantKey, appAccessToken string) (*protocol.GetTenantAccessTokenResp, error) {
	reqData := &protocol.GetTenantAccessTokenISVReq{
		AppAccessToken: appAccessToken,
		TenantKey:      tenantKey,
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.GetTenantAccessTokenIsvPath, nil, reqData)
	if err != nil {
		return nil, fmt.Errorf("doHttpOApiError[%v]", err)
	}

	rspData := &protocol.GetTenantAccessTokenResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes))
	}

	if rspData.Code != 0 {
		return nil, fmt.Errorf("oapiReturnError[code:%d msg:%s]", rspData.Code, rspData.Msg)
	}
	return rspData, nil
}

func getInternalAppAccessToken(ctx context.Context, appID, appSecret string) (*protocol.GetAppAccessTokenInternalResp, error) {
	reqData := &protocol.GetAppAccessTokenInternalReq{
		AppID:     appID,
		AppSecret: appSecret,
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.GetAppAccessTokenInternalPath, nil, reqData)
	if err != nil {
		return nil, fmt.Errorf("doHttpOApiError[%v]", err)
	}

	rspData := &protocol.GetAppAccessTokenInternalResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes))
	}

	if rspData.Code != 0 {
		return nil, fmt.Errorf("oapiReturnError[code:%d msg:%s]", rspData.Code, rspData.Msg)
	}

	return rspData, nil
}

func getIsvAppAccessToken(ctx context.Context, appID, appSecret, appTicket string) (*protocol.GetAppAccessTokenIsvResp, error) {
	reqData := &protocol.GetAppAccessTokenIsvReq{
		AppID:     appID,
		AppSecret: appSecret,
		AppTicket: appTicket,
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.GetAppAccessTokenIsvPath, nil, reqData)
	if err != nil {
		return nil, fmt.Errorf("doHttpOApiError[%v]", err)
	}

	rspData := &protocol.GetAppAccessTokenIsvResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes))
	}

	if rspData.Code != 0 {
		var resultReSend string
		if protocol.ErrAppTicketInvalid == rspData.Code || protocol.ErrAppTicketNil == rspData.Code {
			appTicketManager.SetAppTicket(ctx, appID, "") // disable app ticket

			errReSend := ReSendAppTicket(ctx, appID, appSecret) // appTicket is invalid, resend appTicket
			if errReSend != nil {
				resultReSend = fmt.Sprintf("ReSendAppTicketError[%v]", err)
			} else {
				resultReSend = "ReSendAppTicketSucc"
			}
		} else {
			resultReSend = "ReSendAppTicketDoNothing"
		}

		return nil, fmt.Errorf("oapiReturnError[code:%d msg:%s] %s", rspData.Code, rspData.Msg, resultReSend)
	}

	return rspData, nil
}

func CheckAndDisableTenantToken(ctx context.Context, appID string, tenantKey string, openAPIReturnCode int) {
	if openAPIReturnCode == protocol.ErrTenantAccessTokenInvalid { //common error code

		// openApi return ErrTenantAccessTokenInvalid, need disable local cache
		DisableTenantToken(ctx, appID, tenantKey)
	}
}

func DisableTenantToken(ctx context.Context, appID string, tenantKey string) {
	tokenManager, err := appconfig.GetTokenManager(appID)
	if err != nil {
		common.Logger(ctx).Errorf("SDK-Disable-TenantAccessToken: appID[%s]tenantKey[%s], getManager error[%v]", appID, tenantKey, err)
		return
	}
	tokenManager.DisableTenantAccessToken(tenantKey)

	common.Logger(ctx).Infof("SDK-Disable-TenantAccessToken: appID[%s]tenantKey[%s], disable local cache success", appID, tenantKey)
}

func CheckAndDisableAppToken(ctx context.Context, appID string, openAPIReturnCode int) {
	if openAPIReturnCode == protocol.ErrAppAccessTokenInvalid { //common error code

		// openApi return ErrAppAccessTokenInvalid, need disable local cache
		DisableAppToken(ctx, appID)
	}
}

func DisableAppToken(ctx context.Context, appID string) {
	tokenManager, err := appconfig.GetTokenManager(appID)
	if err != nil {
		common.Logger(ctx).Errorf("SDK-Disable-AppAccessToken: appID[%s], getManager error[%v]", appID, err)
		return
	}
	tokenManager.DisableAppAccessToken()

	common.Logger(ctx).Infof("SDK-Disable-AppAccessToken: appID[%s], disable local cache success", appID)
}
