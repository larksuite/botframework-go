// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

var TplGinCallback = `package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/event"
)

// EventCallback open platform event
func EventCallback(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		common.Logger(c).Errorf("eventReqParamsError: readHttpBodyError err[%v]bodyLen[%d]", err, len(body))
		c.JSON(500, gin.H{"codemsg": common.ErrEventParams.String()})
		return
	}

	appID := "{{.AppID}}"
	challenge, err := event.EventCallback(c, string(body), appID)
	common.Logger(c).Infof("eventInfo: challenge[%s] err[%v]", challenge, err)
	if err != nil {
		c.JSON(500, gin.H{"codemsg": fmt.Sprintf("%v", err)})
	} else if "" != challenge {
		c.JSON(200, gin.H{"challenge": challenge})
	} else {
		c.JSON(200, gin.H{"codemsg": common.Success.String()})
	}
}

// CardCallback card action callback
func CardCallback(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		common.Logger(c).Errorf("eventReqParamsError: readHttpBodyError err[%v]bodyLen[%d]", err, len(body))
		c.JSON(500, gin.H{"codemsg": common.ErrCardParams.String()})
		return
	}

	// for verify signature
	header := map[string]string{
		"X-Lark-Request-Timestamp": c.Request.Header.Get("X-Lark-Request-Timestamp"),
		"X-Lark-Request-Nonce":     c.Request.Header.Get("X-Lark-Request-Nonce"),
		"X-Lark-Signature":         c.Request.Header.Get("X-Lark-Signature"),
	}

	appID := "{{.AppID}}"
	card, challenge, err := event.CardCallBack(c, appID, header, body)
	common.Logger(c).Infof("cardInfo: challenge[%s] err[%v]", challenge, err)
	if err != nil {
		c.JSON(500, gin.H{"codemsg": fmt.Sprintf("%v", err)})
	} else if "" != challenge {
		c.JSON(200, gin.H{"challenge": challenge})
	} else {
		data, err := json.Marshal(card)
		if err != nil {
			c.JSON(500, gin.H{"codemsg": fmt.Sprintf("%v", err)})
		} else {
			c.String(200, string(data))
		}
	}
}
`
