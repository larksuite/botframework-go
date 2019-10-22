// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

const (
	ExpireInterval      = 60 * 5
	ErrAppTicketInvalid = 10012
)

var (
	appTicketManager TicketManager
)

// TicketManager set/get your app-ticket。
type TicketManager interface {
	SetAppTicket(appID, appTicket string) error
	GetAppTicket(appID string) (string, error)
}

func InitISVAppTicketManager(ticketManager TicketManager) error {
	if ticketManager == nil {
		return errors.New("param ticketManager is nil")
	}

	appTicketManager = ticketManager
	return nil
}

func GetTenantAccessToken(ctx context.Context, tenantKey, appID string) (string, error) {
	appToken, err := appconfig.GetToken(appID)
	if err != nil {
		return "", common.ErrAppTokenNotFound.ErrorWithExtErr(err)
	}

	tenantAccessToken, ok := appToken.TenantAccessToken[tenantKey]
	if ok && tenantAccessToken != nil &&
		tenantAccessToken.Token != "" &&
		tenantAccessToken.Expire > time.Now().Unix() {

		return tenantAccessToken.Token, nil
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

	if respBodyObj != nil {
		if appToken.TenantAccessToken[tenantKey] == nil {
			appToken.TenantAccessToken[tenantKey] = new(appconfig.TenantAccessTokenCache)
		}
		appToken.TenantAccessToken[tenantKey].Token = respBodyObj.TenantAccessToken
		appToken.TenantAccessToken[tenantKey].Expire = time.Now().Unix() + int64(respBodyObj.Expire-ExpireInterval)

		return respBodyObj.TenantAccessToken, nil
	}

	return "", common.ErrRespDataIsNil.Error()
}

// ReSendAppTicket app-ticket will be pushed to this service when call this function
func ReSendAppTicket(ctx context.Context, appID, appSecret string) error {
	reqData := &protocol.AppTicketReq{
		AppID:     appID,
		AppSecret: appSecret,
	}

	rspBytes, err := common.DoHttpPostOApi(protocol.ResendAppTicketPath, nil, reqData)
	if err != nil {
		return common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.AppTicketResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return common.ErrJsonUnmarshal.ErrorWithExtErr(err)
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

	err = appTicketManager.SetAppTicket(appTicket.AppID, appTicket.AppTicket)
	if err != nil {
		return common.ErrSetAppTicketFailed.ErrorWithExtErr(err)
	}

	return nil
}

func GetAppAccessToken(ctx context.Context, appID string) (string, error) {
	appToken, err := appconfig.GetToken(appID)
	if err != nil {
		return "", common.ErrAppTokenNotFound.ErrorWithExtErr(err)
	}

	if appToken.AppAccessToken != nil &&
		appToken.AppAccessToken.Token != "" &&
		appToken.AppAccessToken.Expire > time.Now().Unix() {

		return appToken.AppAccessToken.Token, nil
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
		appTicket, err := appTicketManager.GetAppTicket(appID)
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

	if appToken.AppAccessToken == nil {
		appToken.AppAccessToken = new(appconfig.AppAccessTokenCache)
	}
	appToken.AppAccessToken.Token = appAccessToken
	appToken.AppAccessToken.Expire = time.Now().Unix() + int64(expireSecond-ExpireInterval)

	return appAccessToken, nil
}

func getInternalTenantAccessToken(ctx context.Context, appID, appSecret string) (*protocol.GetTenantAccessTokenResp, error) {
	reqData := &protocol.GetTenantAccessTokenInternalReq{
		AppID:     appID,
		AppSecret: appSecret,
	}

	rspBytes, err := common.DoHttpPostOApi(protocol.GetTenantAccessTokenInternalPath, nil, reqData)
	if err != nil {
		return nil, fmt.Errorf("doHttpOApiError[%v]", err)
	}

	rspData := &protocol.GetTenantAccessTokenResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, fmt.Errorf("jsonUnmarshalError[%v]", err)
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

	rspBytes, err := common.DoHttpPostOApi(protocol.GetTenantAccessTokenIsvPath, nil, reqData)
	if err != nil {
		return nil, fmt.Errorf("doHttpOApiError[%v]", err)
	}

	rspData := &protocol.GetTenantAccessTokenResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, fmt.Errorf("jsonUnmarshalError[%v]", err)
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

	rspBytes, err := common.DoHttpPostOApi(protocol.GetAppAccessTokenInternalPath, nil, reqData)
	if err != nil {
		return nil, fmt.Errorf("doHttpOApiError[%v]", err)
	}

	rspData := &protocol.GetAppAccessTokenInternalResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, fmt.Errorf("jsonUnmarshalError[%v]", err)
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

	rspBytes, err := common.DoHttpPostOApi(protocol.GetAppAccessTokenIsvPath, nil, reqData)
	if err != nil {
		return nil, fmt.Errorf("doHttpOApiError[%v]", err)
	}

	rspData := &protocol.GetAppAccessTokenIsvResp{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, fmt.Errorf("jsonUnmarshalError[%v]", err)
	}

	if rspData.Code != 0 {
		var resultReSend string
		if ErrAppTicketInvalid == rspData.Code {
			appTicketManager.SetAppTicket(appID, "")

			errReSend := ReSendAppTicket(ctx, appID, appSecret) // appTicket is invalid, resend appTicket
			if errReSend != nil {
				resultReSend = fmt.Sprintf("ReSendAppTicketError[%v]", err)
			} else {
				resultReSend = "ReSendAppTicketSucc"
			}
		}

		return nil, fmt.Errorf("oapiReturnError[code:%d msg:%s] %s", rspData.Code, rspData.Msg, resultReSend)
	}

	return rspData, nil
}
