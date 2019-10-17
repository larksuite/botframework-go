package auth

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/larksuite/botframework-go/SDK/common"
)

// Default Redis AppTicket Manager: auth.NewDefaultRedisAppTicketManager, need run redis-server
// Default Local AppTicket Manager: auth.NewDefaultLocalAppTicketManager, app-ticket will be lost when service is restarted

// redisAppTicketManager
type redisAppTicketManager struct {
	Client *redis.Client
}

func NewDefaultRedisAppTicketManager(mapParams map[string]string) *redisAppTicketManager {
	r := &redisAppTicketManager{}
	err := r.InitRedis(mapParams)
	if err != nil {
		panic(fmt.Sprintf("New DefaultRedisAppTicketManager Error: InitRedis error[%v]", err))
	}
	return r
}

func (a *redisAppTicketManager) InitRedis(mapParams map[string]string) error {
	a.Client = redis.NewClient(&redis.Options{
		Addr: mapParams["addr"], // addr = host:port , demo: "127.0.0.1:6379"
	})
	_, err := a.Client.Ping().Result()
	if err != nil {
		return fmt.Errorf("init db auth err[%v]", err)
	}

	return nil
}

func (a *redisAppTicketManager) SetAppTicket(appID, appTicket string) error {
	_, err := a.Client.Set("appticket:"+appID, appTicket, 0).Result()
	if err != nil {
		return fmt.Errorf("set auth err[%v]", err)
	}
	return nil
}

func (a *redisAppTicketManager) GetAppTicket(appID string) (string, error) {
	value, err := a.Client.Get("appticket:" + appID).Result()
	if err != nil {
		return "", fmt.Errorf("get auth err[%v]", err)
	}
	return value, nil
}

// localAppTicketManager set/get your app-ticket in-process, app-ticket will be lost when service is restarted
type localAppTicketManager struct {
	appTicketMap map[string]string
}

func NewDefaultLocalAppTicketManager() *localAppTicketManager {
	return &localAppTicketManager{}
}

func (a *localAppTicketManager) SetAppTicket(appID, appTicket string) error {
	if a.appTicketMap == nil {
		a.appTicketMap = make(map[string]string)
	}
	a.appTicketMap[appID] = appTicket
	return nil
}

func (a *localAppTicketManager) GetAppTicket(appID string) (string, error) {
	if appTicket, ok := a.appTicketMap[appID]; ok {
		return appTicket, nil
	}
	return "", common.ErrAppTicketNotFound.Error()
}
