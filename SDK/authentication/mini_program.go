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
func MiniProgramLoginValidate(code string, appID string, appSecret string) (*protocol.MiniProgramLoginResponse, error) {
	// check params
	if code == "" || appID == "" || appSecret == "" {
		return nil, common.ErrValidateParams.ErrorWithExtStr("code/appID/appSecret is empty")
	}

	rspBytes, err := common.DoHttpGetOApi(protocol.MPLoginValidatePath, map[string]string{},
		protocol.GenMiniProgramLoginRequest(code, appID, appSecret))

	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.MiniProgramLoginResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		return nil, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}
