package common

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// Just for the convenience of implementing the default-session-manager/default-appticket-manager

type DefaultRedisClient struct {
	Client *redis.Client
}

func (d *DefaultRedisClient) InitDB(mapParams map[string]string) error {
	d.Client = redis.NewClient(&redis.Options{
		Addr: mapParams["addr"], // addr = host:port , demo: "127.0.0.1:6379"
	})

	_, err := d.Client.Ping().Result()
	if err != nil {
		return fmt.Errorf("init db error[%v]", err)
	}

	return nil
}

func (d *DefaultRedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	if d.Client == nil {
		return fmt.Errorf("db_client isnot initialized, key[%s]", key)
	}

	_, err := d.Client.Set(key, value, expiration).Result()
	if err != nil {
		return fmt.Errorf("set value error[%v], key[%s]", err, key)
	}

	return nil
}

func (d *DefaultRedisClient) Get(key string) (string, error) {
	if d.Client == nil {
		return "", fmt.Errorf("db_client isnot initialized, key[%s]", key)
	}

	value, err := d.Client.Get(key).Result()
	if err != nil {
		return "", fmt.Errorf("get value error[%v], key[%s]", err, key)
	}

	return value, nil
}
