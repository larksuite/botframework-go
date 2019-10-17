// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

var TplRegist = `package handler_event

import (
	"github.com/larksuite/botframework-go/SDK/event"
	{{if .EventList}}"github.com/larksuite/botframework-go/SDK/protocol"{{end}}
)

// If the code is first generated, all code files are generated from the configuration file.
//
// If you modify the configuration file later, and regenerate the code on the original path,
// only the ./handler/regist.go will be forced updated, other files are not updated to avoid overwriting user-defined code.
//
// The ./handler/regist.go file will be forced update, you should not write your business code in the file.

// RegistHandler: regist handler
func RegistHandler(appID string) {

	// regist open platform event handler {{range .EventList}}
	event.EventRegister(appID, protocol.EventType{{.EventName}}, Event{{.EventName}}){{end}}

	// regist bot recv message handler {{range .BotCmdList}}
	event.BotRecvMsgRegister(appID, "{{.Cmd}}", BotRecvMsg{{.FuncName}}){{end}}

	// regist card action handler {{range .CardList}}
	event.CardRegister(appID, "{{.MethodName}}", Action{{.FuncName}}){{end}}
}
`
