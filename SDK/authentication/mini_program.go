// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authentication

import (
	"encoding/json"
	"fmt"

	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

// Code Exchange Token
func MiniProgramValidateByAppToken(code string, appAccessToken string) (*protocol.MiniProgramLoginByAppTokenResponse, error) {
	// check params
	if code == "" || appAccessToken == "" {
		return nil, common.ErrValidateParams.ErrorWithExtStr("code/appAccessToken is empty")
	}

	request := &protocol.MiniProgramLoginByAppTokenRequest{
		Code: code,
	}

	rspBytes, _, err := common.DoHttpPostOApi(protocol.MPValidateByAppTokenPath, common.NewHeaderToken(appAccessToken), request)

	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.MiniProgramLoginByAppTokenResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

// Code Exchange Token
func MiniProgramValidateByIDSecret(code string, appID string, appSecret string) (*protocol.MiniProgramLoginByIDSecretResponse, error) {
	// check params
	if code == "" || appID == "" || appSecret == "" {
		return nil, common.ErrValidateParams.ErrorWithExtStr("code/appID/appSecret is empty")
	}

	rspBytes, _, err := common.DoHttpGetOApi(protocol.MPValidateByIDSecretPath, map[string]string{},
		protocol.GenMiniProgramLoginByIDSecretRequest(code, appID, appSecret))

	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.MiniProgramLoginByIDSecretResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}
