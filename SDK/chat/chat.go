// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

func GetChatInfo(ctx context.Context, tenantKey, appID string, chatID string) (*protocol.GetGroupInfoResponse, error) {
	// check params
	if tenantKey == "" || appID == "" || chatID == "" {
		return nil, common.ErrChatParams.ErrorWithExtStr("param is empty or is nil")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"chat_id": chatID,
	}

	rspBytes, err := common.DoHttpGetOApi(protocol.GetChatInfoPath, common.NewHeaderToken(accessToken), params)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.GetGroupInfoResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		return nil, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}
	return rspData, nil
}

func GetChatList(ctx context.Context, tenantKey, appID string, pageSize int, pageToken string) (*protocol.GetGroupListResponse, error) {
	// check params
	if tenantKey == "" || appID == "" || pageSize <= 0 || pageSize > 200 {
		return nil, common.ErrChatParams.ErrorWithExtStr("param is empty or is invalid")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"page_size":  fmt.Sprintf("%d", pageSize),
		"page_token": pageToken,
	}

	rspBytes, err := common.DoHttpGetOApi(protocol.GetChatListPath, common.NewHeaderToken(accessToken), params)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.GetGroupListResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		return nil, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}
	return rspData, nil
}
