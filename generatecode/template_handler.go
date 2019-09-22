// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

// event template
type EventTemplate struct {
	EventList  []Event      // event
	BotCmdList []BotCommand // bot recv msg
	CardList   []CardAction // bot recv msg
	UseJson    bool
	UseAuth    bool
}

func (tpl *EventTemplate) AddEvent(name string) {
	tpl.EventList = append(tpl.EventList, Event{EventName: name})
}

func (tpl *EventTemplate) AddBotCommand(cmd, description string) {
	tpl.BotCmdList = append(tpl.BotCmdList, BotCommand{Cmd: cmd, FuncName: FormatFuncName(cmd), Description: description})
}

func (tpl *EventTemplate) AddCardAction(name string) {
	tpl.CardList = append(tpl.CardList, CardAction{MethodName: name, FuncName: FormatFuncName(name)})
}

type Event struct {
	EventName string
}

type BotCommand struct {
	Cmd         string
	FuncName    string
	Description string
}
type CardAction struct {
	MethodName string
	FuncName   string
}
