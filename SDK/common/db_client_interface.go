// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import "time"

type DBClient interface {
	InitDB(mapParams map[string]string) error

	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
}
