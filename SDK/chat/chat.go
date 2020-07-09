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
	if appID == "" || chatID == "" {
		return nil, common.ErrChatParams.ErrorWithExtStr("param is empty or is nil")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, statusCode, err := common.DoHttpGetOApi(protocol.GetChatInfoPath, common.NewHeaderToken(accessToken),
		protocol.GenGetGroupInfoRequest(chatID))

	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.GetGroupInfoResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func GetChatList(ctx context.Context, tenantKey, appID string, pageSize int, pageToken string) (*protocol.GetGroupListResponse, error) {
	// check params
	if appID == "" || pageSize <= 0 || pageSize > 200 {
		return nil, common.ErrChatParams.ErrorWithExtStr("param is empty or is invalid")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, statusCode, err := common.DoHttpGetOApi(protocol.GetChatListPath, common.NewHeaderToken(accessToken),
		protocol.GenGetGroupListRequest(pageSize, pageToken))

	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.GetGroupListResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func CreateChat(ctx context.Context, tenantKey, appID string, request *protocol.CreateChatRequest) (*protocol.CreateChatResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrCreateChatParams.ErrorWithExtStr("param is empty or is invalid")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.CreateChatPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.CreateChatResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func UpdateChatInfo(ctx context.Context, tenantKey, appID string, request *protocol.UpdateChatInfoRequest) (*protocol.UpdateChatInfoResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrUpdateChatInfoParams.ErrorWithExtStr("param is empty or is invalid")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.UpdateChatInfoPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.UpdateChatInfoResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func AddUserToChat(ctx context.Context, tenantKey, appID string, request *protocol.AddUserToChatRequest) (*protocol.AddUserToChatResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrAddUserToChatParams.ErrorWithExtStr("param is empty or is invalid")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.AddUserToChatPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.AddUserToChatResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func DeleteUserFromChat(ctx context.Context, tenantKey, appID string, request *protocol.DeleteUserFromChatRequest) (*protocol.DeleteUserFromChatResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrDeleteUserFromChatParams.ErrorWithExtStr("param is empty or is invalid")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.DeleteUserFromChatPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.DeleteUserFromChatResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func DisbandChat(ctx context.Context, tenantKey, appID string, request *protocol.DisbandChatRequest) (*protocol.DisbandChatResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrDisbandChatParams.ErrorWithExtStr("param is empty or is invalid")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, statusCode, err := common.DoHttpPostOApi(protocol.DisbandChatPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.DisbandChatResponse{}
	err = json.Unmarshal(rspBytes, rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(
			fmt.Errorf("jsonUnmarshalError[%v] httpStatusCode[%d] httpBody[%s]", err, statusCode, string(rspBytes)))
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}
