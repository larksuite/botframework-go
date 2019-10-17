// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

var TplEvent = `package handler_event

import (
	"context"
	{{if .UseJson}}"encoding/json"
	"fmt"{{end}}
	{{if .UseAuth}}"github.com/larksuite/botframework-go/SDK/auth"{{end}}
	"github.com/larksuite/botframework-go/SDK/event"
	"github.com/larksuite/botframework-go/SDK/protocol"
)
{{range .EventList}}
// event-{{.EventName}}
func Event{{.EventName}}(ctx context.Context, eventBody []byte) error {
	{{if eq .EventName "Message"}}
	return event.BotRecvMsgHandler(ctx, eventBody){{else if eq .EventName "AppTicket"}}
	return auth.RefreshAppTicket(ctx, eventBody){{else}}
	request := &protocol.{{.EventName}}Event{}
	err := json.Unmarshal(eventBody, request)
	if err != nil {
		return fmt.Errorf("jsonUnmarshalError[%v]", err)
	}

	// NOTE your business code

	return nil{{end}}
}
{{end}}
// handler-bot receive message {{range .BotCmdList}}
// cmd-{{.Cmd}} description: {{.Description}}
func BotRecvMsg{{.FuncName}}(ctx context.Context, msg *protocol.BotRecvMsg) error {
	// NOTE your business code

	return nil
}
{{end}}
`
