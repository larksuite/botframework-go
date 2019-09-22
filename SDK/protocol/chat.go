// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

type GetGroupListRequest struct {
	PageSize  int    `json:"page_size"` // default=100, max=200
	PageToken string `json:"page_token"`
}

type GetGroupListResponse struct {
	BaseResponse
	Data struct {
		Groups    []*Group `json:"groups"`
		HasMore   bool     `json:"has_more"`
		PageToken string   `json:"page_token"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type Group struct {
	Avatar      string `json:"avatar"`
	ChatID      string `json:"chat_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	OwnerOpenID string `json:"owner_open_id"`
	OwnerUserID string `json:"owner_user_id"`
}

type GetGroupInfoRequest struct {
	ChatID string `json:"chat_id"`
}

type GetGroupInfoResponse struct {
	BaseResponse

	Data struct {
		Group

		ChatI18nNames I18nNames    `json:"i18n_names"`
		Members       []UserIDInfo `json:"members"`
	} `json:"data,omitempty" validate:"omitempty"`
}
