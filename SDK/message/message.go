// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package message

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

// SendTextMessage: send text message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  user: the user who you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  text: the text that you will send
func SendTextMessage(ctx context.Context, tenantKey, appID string,
	user *protocol.UserInfo, rootID string,
	text string) (*protocol.SendMsgResponse, error) {

	return sendMsg(ctx, tenantKey, appID, protocol.NewTextMsgReq(user, rootID, text))
}

// SendImageMessage: send image message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  user:  the user who you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  url, path, imageKey: the image that you will send. imageKey > path(The local path name of the image) > url(The URL of the image)
func SendImageMessage(ctx context.Context, tenantKey, appID string,
	user *protocol.UserInfo, rootID string,
	url, path, imageKey string) (*protocol.SendMsgResponse, error) {

	var err error
	if imageKey == "" {
		imageKey, err = GetImageKey(ctx, tenantKey, appID, url, path)
		if err != nil {
			return nil, err
		}
	}

	return sendMsg(ctx, tenantKey, appID, protocol.NewImageMsgReq(user, rootID, imageKey))
}

// SendRichTextMessage: send richtext message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  user:  the user who you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  postForm: the postForm that you will send. You can view the demo code in file richtext_builder_test.go
func SendRichTextMessage(ctx context.Context, tenantKey, appID string,
	user *protocol.UserInfo, rootID string,
	postForm map[protocol.Language]*protocol.RichTextForm) (*protocol.SendMsgResponse, error) {

	// check postForm
	if !checkPostContent(postForm) {
		return nil, common.ErrPostFormParams.Error()
	}

	post := map[string]*protocol.RichTextForm{}
	for k, v := range postForm {
		post[k.String()] = v
	}

	return sendMsg(ctx, tenantKey, appID, protocol.NewPostMsgReq(user, rootID, post))
}

// SendShareChatMessage: send shared chat message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  user:  the user who you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  shareChatID: shared chat id
func SendShareChatMessage(ctx context.Context, tenantKey, appID string,
	user *protocol.UserInfo, rootID string,
	shareChatID string) (*protocol.SendMsgResponse, error) {

	return sendMsg(ctx, tenantKey, appID, protocol.NewShareChatMsgReq(user, rootID, shareChatID))
}

// SendCardMessage: send card message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  user:  the user who you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  card: the cardForm that you will send. You can view the demo code in file card_builder_test.go
// @param  updateMulti: Controls whether the card is a shared card (all users share the same message card), the default is false
func SendCardMessage(ctx context.Context, tenantKey, appID string,
	user *protocol.UserInfo, rootID string,
	card protocol.CardForm, updateMulti bool) (*protocol.SendCardMsgResponse, error) {

	return sendCardMsg(ctx, tenantKey, appID, protocol.NewCardMsgReq(user, rootID, card, updateMulti))
}

// UpdateCard: update card
func UpdateCard(ctx context.Context, tenantKey, appID string, token string, card protocol.CardForm) (*protocol.UpdateCardResponse, error) {
	// check params
	if appID == "" || token == "" {
		return nil, common.ErrCardUpdateParams.ErrorWithExtStr("param is empty or is nil")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	request := &protocol.UpdateCardRequest{
		Token: token,
		Card:  card,
	}

	rspBytes, _, err := common.DoHttpPostOApi(protocol.CardUpdatePath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.UpdateCardResponse{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

// SendTextMessageBatch: batch send text message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  info:  Department id list / user id list, you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  text: the text that you will send
func SendTextMessageBatch(ctx context.Context, tenantKey, appID string,
	info *protocol.BatchBaseInfo, rootID string,
	text string) (*protocol.SendMsgBatchResponse, error) {

	return sendMsgBatch(ctx, tenantKey, appID, protocol.NewBatchTextMsgReq(info, rootID, text))
}

// SendImageMessageBatch: batch send image message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  info:  Department id list / user id list, you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  url, path, imageKey: the image that you will send. imageKey > path(The local path name of the image) > url(The URL of the image)
func SendImageMessageBatch(ctx context.Context, tenantKey, appID string,
	info *protocol.BatchBaseInfo, rootID string,
	url, path, imageKey string) (*protocol.SendMsgBatchResponse, error) {

	var err error
	if imageKey == "" {
		imageKey, err = GetImageKey(ctx, tenantKey, appID, url, path)
		if err != nil {
			return nil, err
		}
	}

	return sendMsgBatch(ctx, tenantKey, appID, protocol.NewBatchImageMsgReq(info, rootID, imageKey))
}

// SendRichTextMessageBatch: batch send richtext message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  info:  Department id list / user id list, you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  postForm: the postForm that you will send. You can view the demo code in file richtext_builder_test.go
func SendRichTextMessageBatch(ctx context.Context, tenantKey, appID string,
	info *protocol.BatchBaseInfo, rootID string,
	postForm map[protocol.Language]*protocol.RichTextForm) (*protocol.SendMsgBatchResponse, error) {

	// check postForm
	if !checkPostContent(postForm) {
		return nil, common.ErrPostFormParams.Error()
	}

	post := map[string]*protocol.RichTextForm{}
	for k, v := range postForm {
		post[k.String()] = v
	}

	return sendMsgBatch(ctx, tenantKey, appID, protocol.NewBatchPostMsgReq(info, rootID, post))
}

// SendShareChatMessageBatch: batch send shared chat message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  info:  Department id list / user id list, you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  shareChatID: shared chat id
func SendShareChatMessageBatch(ctx context.Context, tenantKey, appID string,
	info *protocol.BatchBaseInfo, rootID string,
	shareChatID string) (*protocol.SendMsgBatchResponse, error) {

	return sendMsgBatch(ctx, tenantKey, appID, protocol.NewBatchShareChatMsgReq(info, rootID, shareChatID))
}

// SendCardMessageBatch: batch send card message
// @param  ctx: context
// @param  tenantKey: tenant key. If you don't know it, ask your tenant administrator
// @param  appID: your app ID
// @param  info:  Department id list / user id list, you will send message to
// @param  rootID: The open_message_id of the message that needs to be replied. If do not reply to the message, fill in empty string
// @param  card: the cardForm that you will send. You can view the demo code in file card_builder_test.go
// @param  updateMulti: Controls whether the card is a shared card (all users share the same message card), the default is false
func SendCardMessageBatch(ctx context.Context, tenantKey, appID string,
	info *protocol.BatchBaseInfo, rootID string,
	card protocol.CardForm, updateMulti bool) (*protocol.SendCardMsgBatchResponse, error) {

	return sendCardMsgBatch(ctx, tenantKey, appID, protocol.NewBatchCardMsgReq(info, rootID, card, updateMulti))
}

func sendMsg(ctx context.Context,
	tenantKey, appID string,
	request *protocol.SendMsgRequest) (*protocol.SendMsgResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrSendMsgParams.ErrorWithExtStr("param is empty or is nil")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, _, err := common.DoHttpPostOApi(protocol.SendMessagePath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.SendMsgResponse{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func sendCardMsg(ctx context.Context,
	tenantKey, appID string,
	request *protocol.SendCardMsgRequest) (*protocol.SendCardMsgResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrSendMsgParams.ErrorWithExtStr("param is empty or is nil")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, _, err := common.DoHttpPostOApi(protocol.SendMessagePath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.SendCardMsgResponse{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func sendMsgBatch(ctx context.Context,
	tenantKey, appID string,
	request *protocol.SendMsgBatchRequest) (*protocol.SendMsgBatchResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrSendMsgParams.ErrorWithExtStr("param is empty or is nil")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, _, err := common.DoHttpPostOApi(protocol.SendMessageBatchPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.SendMsgBatchResponse{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}

func sendCardMsgBatch(ctx context.Context,
	tenantKey, appID string,
	request *protocol.SendCardMsgBatchRequest) (*protocol.SendCardMsgBatchResponse, error) {
	// check params
	if appID == "" || request == nil {
		return nil, common.ErrSendMsgParams.ErrorWithExtStr("param is empty or is nil")
	}

	accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
	if err != nil {
		return nil, err
	}

	rspBytes, _, err := common.DoHttpPostOApi(protocol.SendMessageBatchPath, common.NewHeaderToken(accessToken), request)
	if err != nil {
		return nil, common.ErrOpenApiFailed.ErrorWithExtErr(err)
	}

	rspData := &protocol.SendCardMsgBatchResponse{}
	err = json.Unmarshal(rspBytes, &rspData)
	if err != nil {
		return nil, common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}

	if rspData.Code != 0 {
		auth.CheckAndDisableTenantToken(ctx, appID, tenantKey, rspData.Code)
		return rspData, common.ErrOpenApiReturnError.ErrorWithExtStr(fmt.Sprintf("[code:%d msg:%s]", rspData.Code, rspData.Msg))
	}

	return rspData, nil
}
