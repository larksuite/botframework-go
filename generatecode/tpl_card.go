// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

var TplCard = `package handler_event

import (
	"context"
	"encoding/json"

	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)
{{range .CardList}}
// methodName-{{.MethodName}}
func Action{{.FuncName}}(ctx context.Context, callback *protocol.CardCallbackForm) (*protocol.CardForm, error) {
	method, _ := callback.Action.Value["method"]
	sessionID, _ := callback.Action.Value["sid"]
	common.Logger(ctx).Infof("cardActionCallBack: method[%s]sessionID[%s]", method, sessionID)

	// get meta
	meta := &protocol.Meta{}
	if metaData, ok := callback.Action.Value["meta"]; ok {
		_ = json.Unmarshal([]byte(metaData), meta)
	}

	// NOTE your business code

	card := &protocol.CardForm{}

	return card, nil
}
{{end}}
`
