package common

import "time"

type DBClient interface {
	InitDB(mapParams map[string]string) error

	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
}
