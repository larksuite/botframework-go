// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package chat

import (
	"context"
)

func CheckUserIDInGroup(ctx context.Context, tenantKey, appID, chatID, userID string) (bool, error) {
	rspData, err := GetChatInfo(ctx, tenantKey, appID, chatID)
	if err != nil {
		return false, err
	}

	for _, v := range rspData.Data.Members {
		if userID == v.UserID {
			return true, nil
		}
	}

	return false, nil
}

func CheckOpenIDInGroup(ctx context.Context, tenantKey, appID, chatID, openID string) (bool, error) {
	rspData, err := GetChatInfo(ctx, tenantKey, appID, chatID)
	if err != nil {
		return false, err
	}

	for _, v := range rspData.Data.Members {
		if openID == v.OpenID {
			return true, nil
		}
	}

	return false, nil
}

func CheckBotInGroup(ctx context.Context, tenantKey, appID, chatID string) (bool, error) {
	const MaxGetPage int = 1000 //  avoid to fall into an endless loop

	pageToken := ""
	for page := 1; page < MaxGetPage; page++ {
		rspData, err := GetChatList(ctx, tenantKey, appID, 100, pageToken)
		if err != nil {
			return false, err
		}

		for _, v := range rspData.Data.Groups {
			if chatID == v.ChatID {
				return true, nil
			}
		}

		if !rspData.Data.HasMore {
			return false, nil
		}

		pageToken = rspData.Data.PageToken
	}

	return false, nil
}

func CheckUserIDBotInSameGroup(ctx context.Context, tenantKey, appID, chatID, userID string) (bool, error) {
	isBotInGroup, err := CheckBotInGroup(ctx, tenantKey, appID, chatID)
	if err != nil {
		return false, err
	}
	isUserIDInGroup, err := CheckUserIDInGroup(ctx, tenantKey, appID, chatID, userID)
	if err != nil {
		return false, err
	}

	if isBotInGroup && isUserIDInGroup {
		return true, nil
	}
	return false, nil
}

func CheckOpenIDBotInSameGroup(ctx context.Context, tenantKey, appID, chatID, openID string) (bool, error) {
	isBotInGroup, err := CheckBotInGroup(ctx, tenantKey, appID, chatID)
	if err != nil {
		return false, err
	}
	isOpenIDInGroup, err := CheckOpenIDInGroup(ctx, tenantKey, appID, chatID, openID)
	if err != nil {
		return false, err
	}

	if isBotInGroup && isOpenIDInGroup {
		return true, nil
	}
	return false, nil
}
