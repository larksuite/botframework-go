// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

type HandlerBotMsg func(ctx context.Context, msg *protocol.BotRecvMsg) error

// CommandHandlerManager cmd --> handler
type CommandHandlerManager struct {
	mapHandler map[string]map[string]HandlerBotMsg
}

func (p *CommandHandlerManager) Set(appID string, cmdName string, handler HandlerBotMsg) {
	if handler == nil {
		return
	}
	if _, ok := p.mapHandler[appID]; !ok {
		p.mapHandler[appID] = make(map[string]HandlerBotMsg, 0)
	}

	cmdName = strings.ToLower(cmdName)
	p.mapHandler[appID][cmdName] = handler
}

func (p *CommandHandlerManager) Get(appID string, cmdName string) (HandlerBotMsg, error) {
	cmdName = strings.ToLower(cmdName)

	if _, ok := p.mapHandler[appID]; !ok {
		return nil, fmt.Errorf("botRecvMsg appid[%s] has not been registered", appID)
	}
	if _, ok := p.mapHandler[appID][cmdName]; !ok {
		return nil, fmt.Errorf("botRecvMsg cmdName[%s] has not been registered in appid[%s]", cmdName, appID)
	}

	return p.mapHandler[appID][cmdName], nil
}

var cmdHandler *CommandHandlerManager

func init() {
	cmdHandler = &CommandHandlerManager{mapHandler: make(map[string]map[string]HandlerBotMsg, 0)}
}

// BotRecvMsgRegister appid+cmd --> handler
func BotRecvMsgRegister(appID string, cmdName string, handler HandlerBotMsg) error {
	if appID == "" || cmdName == "" {
		return common.ErrBotRecvMsgRegister.ErrorWithExtStr("appID/cmdName is empty")
	}
	if handler == nil {
		return common.ErrBotRecvMsgRegister.ErrorWithExtStr("action handler is nil")
	}

	cmdHandler.Set(appID, cmdName, handler)
	return nil
}

// BotRecvMsgHandler callback botRecvMsg
func BotRecvMsgHandler(ctx context.Context, data []byte) error {
	// get msg_type and app_id
	jsonMsg, err := simplejson.NewJson(data)
	if err != nil {
		return common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}
	msgType, err := jsonMsg.Get("msg_type").String()
	if err != nil {
		return common.ErrBotRecvMsgMsgTypeJson.ErrorWithExtErr(err)
	}
	appID, err := jsonMsg.Get("app_id").String()
	if err != nil || appID == "" {
		return common.ErrBotRecvMsgAppIDJson.ErrorWithExtErr(err)
	}

	var msg protocol.BotRecvMsg
	msg.AppID = appID
	msg.MsgType = msgType

	var textWithoutAtBot string
	switch msgType {
	case protocol.EventMsgTypeText:
		msgEvent := &protocol.TextMsgEvent{}
		err := json.Unmarshal(data, msgEvent)
		if err != nil {
			return common.ErrJsonUnmarshal.ErrorWithExtErr(err)
		}
		msg.TenantKey = msgEvent.TenantKey
		msg.RootID = msgEvent.RootID
		msg.ParentID = msgEvent.ParentID
		msg.OpenChatID = msgEvent.OpenChatID
		msg.ChatType = msgEvent.ChatType
		msg.OpenID = msgEvent.OpenID
		msg.OpenMessageID = msgEvent.OpenMessageID
		msg.OriData = msgEvent

		textWithoutAtBot = msgEvent.TextWithoutAtBot
	case protocol.EventMsgTypePost:
		msgEvent := &protocol.PostMsgEvent{}
		err := json.Unmarshal(data, msgEvent)
		if err != nil {
			return common.ErrJsonUnmarshal.ErrorWithExtErr(err)
		}
		msg.TenantKey = msgEvent.TenantKey
		msg.RootID = msgEvent.RootID
		msg.ParentID = msgEvent.ParentID
		msg.OpenChatID = msgEvent.OpenChatID
		msg.ChatType = msgEvent.ChatType
		msg.OpenID = msgEvent.OpenID
		msg.OpenMessageID = msgEvent.OpenMessageID
		msg.OriData = msgEvent

		textWithoutAtBot = ""
	case protocol.EventMsgTypeImage:
		msgEvent := &protocol.ImageMsgEvent{}
		err := json.Unmarshal(data, msgEvent)
		if err != nil {
			return common.ErrJsonUnmarshal.ErrorWithExtErr(err)
		}
		msg.TenantKey = msgEvent.TenantKey
		msg.RootID = msgEvent.RootID
		msg.ParentID = msgEvent.ParentID
		msg.OpenChatID = msgEvent.OpenChatID
		msg.ChatType = msgEvent.ChatType
		msg.OpenID = msgEvent.OpenID
		msg.OpenMessageID = msgEvent.OpenMessageID
		msg.OriData = msgEvent

		textWithoutAtBot = ""
	case protocol.EventMsgTypeMergeForward:
		msgEvent := &protocol.MergeForwardMsgEvent{}
		err := json.Unmarshal(data, msgEvent)
		if err != nil {
			return common.ErrJsonUnmarshal.ErrorWithExtErr(err)
		}
		msg.TenantKey = msgEvent.TenantKey
		msg.RootID = msgEvent.RootID
		msg.ParentID = msgEvent.ParentID
		msg.OpenChatID = msgEvent.OpenChatID
		msg.ChatType = msgEvent.ChatType
		msg.OpenID = msgEvent.OpenID
		msg.OpenMessageID = msgEvent.OpenMessageID
		msg.OriData = msgEvent

		textWithoutAtBot = ""
	default:
		textWithoutAtBot = ""
	}

	msg.TextParam = textWithoutAtBot

	//get cmd
	var cmd string
	s := strings.Split(strings.Trim(msg.TextParam, " "), " ")
	if len(s) > 1 {
		cmd = strings.ToLower(s[0])
		msg.TextParam = strings.Trim(strings.Join(s[1:], " "), " ")
	} else if len(s) > 0 {
		cmd = strings.ToLower(s[0])
		msg.TextParam = ""
	} else {
		cmd = protocol.CmdDefault
	}

	//get handler
	handler, err := cmdHandler.Get(appID, cmd)
	if err != nil {
		if protocol.CmdDefault != cmd {
			cmd = protocol.CmdDefault
			msg.TextParam = textWithoutAtBot

			handler, err = cmdHandler.Get(appID, cmd)
		}

		if err != nil {
			return common.ErrBotRecvMsgHandlerNoFound.ErrorWithExtErr(err)
		}
	}

	err = handler(ctx, &msg)
	if err != nil {
		return common.ErrBotRecvMsgHandlerFailed.ErrorWithExtErr(err)
	}

	return nil
}
