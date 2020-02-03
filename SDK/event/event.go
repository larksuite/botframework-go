// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

type EventHandler func(ctx context.Context, eventBody []byte) error

type EventHandlerManager struct {
	mapHandler map[string]map[string]EventHandler //[app_id][eventType][function]
}

func (a *EventHandlerManager) Set(appID, eventTypeList string, eventHandler EventHandler) {
	if eventHandler == nil {
		return
	}

	if m, ok := a.mapHandler[appID]; !ok || m == nil {
		a.mapHandler[appID] = make(map[string]EventHandler, 0)
	}

	s := strings.Split(strings.Trim(eventTypeList, " "), ",")
	for _, eventType := range s {
		a.mapHandler[appID][eventType] = eventHandler
	}
}

func (a *EventHandlerManager) Get(appID string, eventType string) (EventHandler, error) {
	if _, ok := a.mapHandler[appID]; !ok {
		return nil, common.ErrEventAppIDUnregistered.ErrorWithExtStr(fmt.Sprintf("appid[%s]eventType[%s]", appID, eventType))
	}
	if _, ok := a.mapHandler[appID][eventType]; !ok {
		return nil, common.ErrEventTypeUnregistered.ErrorWithExtStr(fmt.Sprintf("appid[%s]eventType[%s]", appID, eventType))
	}
	if a.mapHandler[appID][eventType] == nil {
		return nil, common.ErrEventHandlerIsNil.ErrorWithExtStr(fmt.Sprintf("appid[%s]eventType[%s]", appID, eventType))
	}

	return a.mapHandler[appID][eventType], nil
}

var eventManager *EventHandlerManager

func init() {
	eventManager = &EventHandlerManager{
		mapHandler: make(map[string]map[string]EventHandler, 0),
	}
}

func EventRegister(appID, eventTypeList string, eventHandler EventHandler) error {
	// 参数校验
	if appID == "" || eventTypeList == "" || eventHandler == nil {
		return common.ErrEventTypeRegister.ErrorWithExtStr(
			fmt.Sprintf("params is empty or nil. AppID[%s]EventType[%s]HandlerIsNil[%t]", appID, eventTypeList, eventHandler == nil))
	}

	eventManager.Set(appID, eventTypeList, eventHandler)

	return nil
}

func EventCallback(ctx context.Context, body string, appID string) (string, error) {
	// check params
	if body == "" || appID == "" {
		return "", common.ErrEventParams.Error()
	}

	// get app config
	appConf, err := appconfig.GetConfig(appID)
	if err != nil {
		return "", common.ErrAppConfNotFound.ErrorWithExtErr(err)
	}

	// decrypt data
	var content string
	if appConf.EncryptKey != "" {
		content, err = eventDataDecrypter(body, appConf.EncryptKey)
		if err != nil {
			return "", common.ErrEventDecrypt.ErrorWithExtErr(err)
		}
	} else {
		content = body
	}

	var callbackBase protocol.CallbackBase
	err = json.Unmarshal([]byte(content), &callbackBase)
	if err != nil {
		return "", common.ErrEventGetBase.ErrorWithExtErr(err)
	}

	// check token
	if callbackBase.Token != appConf.VerifyToken {
		return "", common.ErrEventVeriToken.ErrorWithExtStr(
			fmt.Sprintf("eventVTokenSize[%d]confVTokenSize[%d]", len(callbackBase.Token), len(appConf.VerifyToken)))
	}

	// dispatch event type
	switch callbackBase.Type {
	case protocol.EventChallenge:
		if appConf.AppType == protocol.ISVApp {
			auth.ReSendAppTicket(ctx, appConf.AppID, appConf.AppSecret)
		}

		return callbackBase.Challenge, nil
	case protocol.EventCallback:
		err = eventCallbackHandler(ctx, appID, content)
		if err != nil {
			return "", err
		}
	default:
		return "", common.ErrEventTypeUnknown.ErrorWithExtStr(fmt.Sprintf("appid[%s]type[%s]", appID, callbackBase.Type))
	}

	return "", nil
}

func eventCallbackHandler(ctx context.Context, appID, content string) error {
	// get type and app_id
	jsonBody, err := simplejson.NewJson([]byte(content))
	if err != nil {
		return common.ErrJsonUnmarshal.ErrorWithExtErr(err)
	}
	jsonEvent := jsonBody.Get("event")
	if jsonEvent.Interface() == nil {
		return common.ErrEventGetJsonEvent.Error()
	}
	eventType, err := jsonEvent.Get("type").String()
	if err != nil {
		return common.ErrEventGetJsonType.ErrorWithExtErr(err)
	}
	eventAppID, err := jsonEvent.Get("app_id").String()
	if err != nil {
		return common.ErrEventGetJsonAppID.ErrorWithExtErr(err)
	}

	// check params
	if eventAppID != appID {
		return common.ErrEventAppIDNotMatch.Error()
	}

	// dispatch event type
	handler, err := eventManager.Get(appID, eventType)
	if err != nil {
		return err
	}

	byteEvent, err := jsonEvent.MarshalJSON()
	if err != nil {
		return common.ErrJsonMarshal.ErrorWithExtErr(err)
	}

	err = handler(ctx, byteEvent)
	if err != nil {
		return common.ErrEventHandlerFailed.ErrorWithExtErr(err)
	}

	return nil
}

func eventDataDecrypter(encryptData, keyStr string) (string, error) {
	type AESMsg struct {
		Encrypt string `json:"encrypt"`
	}
	var encrypt AESMsg
	err := json.Unmarshal([]byte(encryptData), &encrypt)
	if err != nil {
		return "", fmt.Errorf("dataDecrypter jsonUnmarshalError[%v]", err)
	}

	buf, err := base64.StdEncoding.DecodeString(encrypt.Encrypt)
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}

	key := sha256.Sum256([]byte(keyStr))

	block, err := aes.NewCipher(key[:sha256.Size])
	if err != nil {
		return "", fmt.Errorf("AESNewCipher Error[%v]", err)
	}

	if len(buf) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)

	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}

	return string(buf[n : m+1]), nil
}
