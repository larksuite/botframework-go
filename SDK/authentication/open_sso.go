// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

// generate open sso authentication URL
func OpenSSOGenerateAuthURL(redirectURL string, appID string, state string) string {
	m := make(url.Values)
	m.Set("redirect_uri", redirectURL)
	m.Set("app_id", appID)
	m.Set("state", state)

	reqURL := common.GetOpenPlatformHost() + string(protocol.OpenSSOGetCodePath) + "?" + m.Encode()

	return reqURL
}

func OpenSSOCodeValidateByAppToken(ctx context.Context, code string, appAccessToken string) (*protocol.OpenSSOTokenResponse, error) {
	tokenReq := &protocol.OpenSSOTokenRequest{
		AppAccessToken: appAccessToken,
		GrantType:      protocol.GrantTypeAuthCode,
		Code:           code,
	}

	return openSSOValidate(ctx, tokenReq)
}

func OpenSSOCodeValidateByIDSecret(ctx context.Context, code string, appID, appSecret string) (*protocol.OpenSSOTokenResponse, error) {
	tokenReq := &protocol.OpenSSOTokenRequest{
		AppID:     appID,
		AppSecret: appSecret,
		GrantType: protocol.GrantTypeAuthCode,
		Code:      code,
	}

	return openSSOValidate(ctx, tokenReq)
}

func OpenSSORefreshTokenByAppToken(ctx context.Context, refreshToken string, appAccessToken string) (*protocol.OpenSSOTokenResponse, error) {
	tokenReq := &protocol.OpenSSOTokenRequest{
		AppAccessToken: appAccessToken,
		GrantType:      protocol.GrantTypeRefreshToken,
		RefreshToken:   refreshToken,
	}

	return openSSOValidate(ctx, tokenReq)
}

func OpenSSORefreshTokenByIDSecret(ctx context.Context, refreshToken string, appID, appSecret string) (*protocol.OpenSSOTokenResponse, error) {
	tokenReq := &protocol.OpenSSOTokenRequest{
		AppID:        appID,
		AppSecret:    appSecret,
		GrantType:    protocol.GrantTypeRefreshToken,
		RefreshToken: refreshToken,
	}

	return openSSOValidate(ctx, tokenReq)
}

func openSSOValidate(ctx context.Context, tokenReq *protocol.OpenSSOTokenRequest) (*protocol.OpenSSOTokenResponse, error) {
	// check params
	if (tokenReq.AppID == "" || tokenReq.AppSecret == "") && tokenReq.AppAccessToken == "" {
		return nil, common.ErrValidateParams.ErrorWithExtStr("(app_id+app_secret)/app_access_token are empty")
	}
	if tokenReq.GrantType != protocol.GrantTypeAuthCode && tokenReq.GrantType != protocol.GrantTypeRefreshToken {
		return nil, common.ErrValidateParams.ErrorWithExtStr("grant_type is invalid")
	}
	if tokenReq.GrantType == protocol.GrantTypeAuthCode && tokenReq.Code == "" {
		return nil, common.ErrValidateParams.ErrorWithExtStr("code is empty")
	}
	if tokenReq.GrantType == protocol.GrantTypeRefreshToken && tokenReq.RefreshToken == "" {
		return nil, common.ErrValidateParams.ErrorWithExtStr("refresh_token is empty")
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.OpenSSOValidatePath, common.NewHeaderJson(), tokenReq)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.OpenSSOTokenResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}
