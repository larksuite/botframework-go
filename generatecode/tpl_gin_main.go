// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

var TplGinMain = `package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
	"{{.Path}}/handler_event"
)

func main() {
	r := gin.Default()

	common.InitLogger(common.NewCommonLogger(), common.DefaultOption())
	defer common.FlushLogger()

	err := InitInfo()
	if err != nil {
		common.Logger(context.TODO()).Errorf("InitError[%v]", err)
		return
	}

	r.POST("{{.EventWebhook}}", EventCallback) //open platform event callback
	r.POST("{{.CardWebhook}}", CardCallback)   //card action callback

	// NOTE your business code

	r.Run(":8089")
}

func InitInfo() error {
	// Initialize app config
	conf := appconfig.AppConfig{
		AppID:   "{{.AppID}}",
		AppType: {{if .IsISVApp}}protocol.ISVApp{{else}}protocol.InternalApp{{end}}, // Independent Software Vendor App / Internal App

		// NOTE your business code
		// get appinfo(app_secret、veri_token、encrypt_key) from redis/mysql or remote config system
		// redis/mysql or remote config system is recommended

		// AppSecret:   redis.GetString("{{.AppID}}" + "AppSecret"),
		// VerifyToken: redis.GetString("{{.AppID}}" + "VerifyToken"),
		// EncryptKey:  redis.GetString("{{.AppID}}" + "EncryptKey"),
	}

	appconfig.Init(conf)

	// ISVApp Set TicketManager
	if conf.AppType == protocol.ISVApp {
		// ISVApp need to implement the TicketManager interface
		// It is recommended to set/get your app-ticket in redis
		// Default Redis AppTicket Manager: auth.NewDefaultRedisAppTicketManager, need run redis-server
		err := auth.InitISVAppTicketManager(auth.NewDefaultRedisAppTicketManager(map[string]string{"addr": "127.0.0.1:6379"}))
		if err != nil {
			return fmt.Errorf("Authorization Initialize Error[%v]", err)
		}
	}

	// regist handler
	handler_event.RegistHandler(conf.AppID)

	return nil
}
`
