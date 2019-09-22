// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
	"github.com/larksuite/botframework-go/generatecode"
	"github.com/jinzhu/configor"
)

var (
	help     bool
	filePath string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&filePath, "f", "", "config file path")
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	if filePath == "" {
		flag.Usage()
		return
	}

	common.InitLogger(common.DefaultOption())
	defer common.FlushLogger()
	ctx := context.Background()

	// read config
	config := &generatecode.GenCodeInfo{}
	err := configor.Load(config, filePath)
	if err != nil {
		common.Logger(ctx).Errorf("genConfigError error[%v]confPath[%s]", err, filePath)
		return
	}

	// check config
	if config.ServiceInfo.Path == "" ||
		config.ServiceInfo.EventWebhook == "" ||
		config.ServiceInfo.CardWebhook == "" ||
		config.ServiceInfo.AppID == "" {

		common.Logger(ctx).Errorf("config error, Path/EventWebhook/CardWebhookAppID can not be empty ")
		return
	}

	// generate code
	GenCodeGin(ctx, config)
}

func GenCodeGin(ctx context.Context, config *generatecode.GenCodeInfo) {
	mainTpl := generatecode.MainTemplate{
		Path:         config.ServiceInfo.Path,
		GenCodePath:  config.ServiceInfo.GenCodePath,
		EventWebhook: config.ServiceInfo.EventWebhook,
		CardWebhook:  config.ServiceInfo.CardWebhook,
		AppID:        config.ServiceInfo.AppID,
		IsISVApp:     config.ServiceInfo.IsISVApp,
	}

	var path string
	if mainTpl.GenCodePath != "" {
		path = mainTpl.GenCodePath
	} else {
		GOPATH := os.Getenv("GOPATH")
		path = GOPATH + "/src/" + mainTpl.Path
	}

	// init path
	err := generatecode.InitPath(path)
	if err != nil {
		common.Logger(ctx).Errorf("initPathError[%v]", err)
		return
	}

	common.Logger(ctx).Infof("generate Gin code in path[%s]", path)
	// generatecode main、callback
	err = generatecode.GenerateCode(ctx, "tplmain", generatecode.TplGinMain, path, "/main.go", mainTpl, false)
	if err != nil {
		common.Logger(ctx).Errorf("generateCodeError[%v]", err)
		return
	}

	err = generatecode.GenerateCode(ctx, "tplcallback", generatecode.TplGinCallback, path, "/callback.go", mainTpl, false)
	if err != nil {
		common.Logger(ctx).Errorf("generateCodeError[%v]", err)
		return
	}

	// generatecode handler
	eventTpl := &generatecode.EventTemplate{}
	for _, v := range config.EventList {
		if v.EventName == "AppTicket" {
			eventTpl.UseAuth = true
		} else if v.EventName == "Message" {

		} else {
			eventTpl.UseJson = true
		}
		eventTpl.AddEvent(v.EventName)
	}

	isHasDefault := false
	for _, v := range config.CommandList {
		if protocol.CmdDefault == strings.ToLower(v.Cmd) {
			isHasDefault = true
		}
		eventTpl.AddBotCommand(v.Cmd, v.Description)
	}
	if !isHasDefault {
		eventTpl.AddBotCommand(protocol.CmdDefault, protocol.DescDefault)
	}

	for _, v := range config.CardActionList {
		eventTpl.AddCardAction(v.MethodName)
	}

	err = generatecode.GenerateCode(ctx, "tplregist", generatecode.TplRegist, path, "/handler/regist.go", eventTpl, true)
	if err != nil {
		common.Logger(ctx).Errorf("generateCodeError[%v]", err)
		return
	}

	err = generatecode.GenerateCode(ctx, "tplevent", generatecode.TplEvent, path, "/handler/event.go", eventTpl, false)
	if err != nil {
		common.Logger(ctx).Errorf("generateCodeError[%v]", err)
		return
	}

	err = generatecode.GenerateCode(ctx, "tplCard", generatecode.TplCard, path, "/handler/card.go", eventTpl, false)
	if err != nil {
		common.Logger(ctx).Errorf("generateCodeError[%v]", err)
		return
	}

	common.Logger(ctx).Infof("Success generateCodeSucc service in path[%s]", path)
	return
}
