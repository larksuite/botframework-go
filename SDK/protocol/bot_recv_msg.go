// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

const (
	CmdDefault  = "default"
	DescDefault = "defalut cmd"
)

type BotRecvMsg struct {
	AppID         string
	TenantKey     string
	MsgType       string
	TextParam     string
	RootID        string
	ParentID      string
	OpenChatID    string
	ChatType      string
	OpenID        string
	OpenMessageID string
	OriData       interface{}
}
