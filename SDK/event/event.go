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

var (
	appID2TypeHandler = make(map[string]map[string]EventHandler) // [app_id][eventType][function]
)

func EventRegister(appID, eventType string, eventHandler EventHandler) error {
	// 参数校验
	if appID == "" {
		return common.ErrEventTypeRegister.ErrorWithExtStr("app id is empty")
	}
	if eventType == "" {
		return common.ErrEventTypeRegister.ErrorWithExtStr("event type is empty")
	}
	if eventHandler == nil {
		return common.ErrEventTypeRegister.ErrorWithExtStr("event handler is nil")
	}

	if appID2TypeHandler == nil {
		return common.ErrEventManagerNotInit.Error()
	}

	if typeHandler, ok := appID2TypeHandler[appID]; ok {
		if typeHandler == nil {
			typeHandler = make(map[string]EventHandler)
			appID2TypeHandler[appID] = typeHandler
		}

		s := strings.Split(strings.Trim(eventType, " "), ",")
		for _, v := range s {
			typeHandler[v] = eventHandler
		}

	} else {
		typeHandler = make(map[string]EventHandler)
		appID2TypeHandler[appID] = typeHandler

		s := strings.Split(strings.Trim(eventType, " "), ",")
		for _, v := range s {
			typeHandler[v] = eventHandler
		}
	}

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
		return "", common.ErrEventVeriToken.Error()
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
	var eventHandler map[string]EventHandler
	var handler EventHandler
	var ok bool

	if eventHandler, ok = appID2TypeHandler[appID]; !ok {
		return common.ErrEventAppIDUnregistered.Error()
	}
	if handler, ok = eventHandler[eventType]; !ok {
		return common.ErrEventTypeUnregistered.Error()
	}
	if handler == nil {
		return common.ErrEventHandlerIsNil.Error()
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
