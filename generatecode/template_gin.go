// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

type MainTemplate struct {
	Path         string
	GenCodePath  string
	EventWebhook string
	CardWebhook  string
	AppID        string
	IsISVApp     bool
}
