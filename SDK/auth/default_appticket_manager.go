// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"context"

	"github.com/larksuite/botframework-go/SDK/common"
)

// DefaultAppTicketManager
type DefaultAppTicketManager struct {
	Client common.DBClient
}

// NewDefaultAppTicketManager demo:
// client := &common.DefaultRedisClient{}
// err := client.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
// if err != nil {
// 	return fmt.Errorf("init redis error[%v]", err)
// }
// manager := auth.NewDefaultAppTicketManager(client)
func NewDefaultAppTicketManager(client common.DBClient) *DefaultAppTicketManager {
	r := &DefaultAppTicketManager{
		Client: client,
	}
	return r
}

func (a *DefaultAppTicketManager) SetAppTicket(ctx context.Context, appID, appTicket string) error {
	return a.Client.Set("appticket:"+appID, appTicket, 0)
}

func (a *DefaultAppTicketManager) GetAppTicket(ctx context.Context, appID string) (string, error) {
	return a.Client.Get("appticket:" + appID)
}
