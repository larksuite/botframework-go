package auth

import (
	"github.com/larksuite/botframework-go/SDK/common"
)

// DefaultAppTicketManager
type DefaultAppTicketManager struct {
	Client common.DBClient
}

// NewDefaultAppTicketManager demo:
// client := &common.DefaultRedisClient{}
// client.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
// manager := auth.NewDefaultAppTicketManager(client)
func NewDefaultAppTicketManager(client common.DBClient) *DefaultAppTicketManager {
	r := &DefaultAppTicketManager{
		Client: client,
	}
	return r
}

func (a *DefaultAppTicketManager) SetAppTicket(appID, appTicket string) error {
	return a.Client.Set("appticket:"+appID, appTicket, 0)
}

func (a *DefaultAppTicketManager) GetAppTicket(appID string) (string, error) {
	return a.Client.Get("appticket:" + appID)
}
