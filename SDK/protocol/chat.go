// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

import "fmt"

func GenGetGroupListRequest(pageSize int, pageToken string) map[string]string {
	return map[string]string{
		"page_size":  fmt.Sprintf("%d", pageSize),
		"page_token": pageToken,
	}
}

type GetGroupListResponse struct {
	BaseResponse
	Data struct {
		Groups    []*Group `json:"groups,omitempty" validate:"omitempty"`
		HasMore   bool     `json:"has_more,omitempty" validate:"omitempty"`
		PageToken string   `json:"page_token,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type Group struct {
	Avatar      string `json:"avatar,omitempty" validate:"omitempty"`
	ChatID      string `json:"chat_id,omitempty" validate:"omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Name        string `json:"name,omitempty" validate:"omitempty"`
	OwnerOpenID string `json:"owner_open_id,omitempty" validate:"omitempty"`
	OwnerUserID string `json:"owner_user_id,omitempty" validate:"omitempty"`
}

func GenGetGroupInfoRequest(chatID string) map[string]string {
	return map[string]string{
		"chat_id": chatID,
	}
}

type GetGroupInfoResponse struct {
	BaseResponse

	Data struct {
		Group

		ChatI18nNames I18nNames    `json:"i18n_names,omitempty" validate:"omitempty"`
		Members       []UserIDInfo `json:"members,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type UpdateChatInfoRequest struct {
	ChatID        string     `json:"chat_id,omitempty" validate:"omitempty"`
	OwnerUserID   *string    `json:"owner_user_id,omitempty" validate:"omitempty"`
	OwnerOpenID   *string    `json:"owner_open_id,omitempty" validate:"omitempty"`
	Name          *string    `json:"name,omitempty" validate:"omitempty"`
	ChatI18nNames *I18nNames `json:"i18n_names,omitempty" validate:"omitempty"`
}

type UpdateChatInfoResponse struct {
	BaseResponse
	Data struct {
		ChatID string `json:"chat_id,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type CreateChatRequest struct {
	Name          string     `json:"name,omitempty" validate:"omitempty"`
	Description   string     `json:"description,omitempty" validate:"omitempty"`
	UserIDs       []string   `json:"user_ids,omitempty" validate:"omitempty"`
	OpenIDs       []string   `json:"open_ids,omitempty" validate:"omitempty"`
	ChatI18nNames *I18nNames `json:"i18n_names,omitempty" validate:"omitempty"`
}

type CreateChatResponse struct {
	BaseResponse
	Data struct {
		ChatID         string   `json:"chat_id,omitempty" validate:"omitempty"`
		InvalidOpenIDs []string `json:"invalid_open_ids,omitempty" validate:"omitempty"`
		InvalidUserIDs []string `json:"invalid_user_ids,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type AddUserToChatRequest struct {
	ChatID  string   `json:"chat_id,omitempty" validate:"omitempty"`
	UserIDs []string `json:"user_ids,omitempty" validate:"omitempty"`
	OpenIDs []string `json:"open_ids,omitempty" validate:"omitempty"`
}

type AddUserToChatResponse struct {
	BaseResponse
	Data struct {
		InvalidOpenIDs []string `json:"invalid_open_ids,omitempty" validate:"omitempty"`
		InvalidUserIDs []string `json:"invalid_user_ids,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type DeleteUserFromChatRequest struct {
	ChatID  string   `json:"chat_id,omitempty" validate:"omitempty"`
	UserIDs []string `json:"user_ids,omitempty" validate:"omitempty"`
	OpenIDs []string `json:"open_ids,omitempty" validate:"omitempty"`
}

type DeleteUserFromChatResponse struct {
	BaseResponse
	Data struct {
		InvalidOpenIDs []string `json:"invalid_open_ids,omitempty" validate:"omitempty"`
		InvalidUserIDs []string `json:"invalid_user_ids,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type DisbandChatRequest struct {
	ChatID string `json:"chat_id,omitempty" validate:"omitempty"`
}

type DisbandChatResponse struct {
	BaseResponse
}
