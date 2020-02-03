// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

type ActionMethod func(ctx context.Context, cardCallback *protocol.CardCallbackForm) (*protocol.CardForm, error)

type ActionHandlerManager struct {
	mapHandler map[string]map[string]ActionMethod // appID => method => func
	ignoreSign map[string]bool
}

func (a *ActionHandlerManager) Set(appID string, method string, v ActionMethod) {
	if v == nil {
		return
	}

	if m, ok := a.mapHandler[appID]; !ok || m == nil {
		a.mapHandler[appID] = make(map[string]ActionMethod, 0)
	}

	a.mapHandler[appID][method] = v
}

func (a *ActionHandlerManager) Get(appID string, method string) (ActionMethod, error) {
	if _, ok := a.mapHandler[appID]; !ok {
		return nil, fmt.Errorf("getCardActionHandler appid[%s] has not been registered", appID)
	}
	if _, ok := a.mapHandler[appID][method]; !ok {
		return nil, fmt.Errorf("getCardActionHandler method[%s] has not been registered in appid[%s]", method, appID)
	}

	return a.mapHandler[appID][method], nil
}

var cardHandler *ActionHandlerManager

func init() {
	cardHandler = &ActionHandlerManager{
		mapHandler: make(map[string]map[string]ActionMethod, 0),
		ignoreSign: make(map[string]bool, 0),
	}
}

func IgnoreSign(appid string, ignore bool) {
	cardHandler.ignoreSign[appid] = ignore
}

func CardRegister(appID string, method string, handler ActionMethod) error {
	if appID == "" || method == "" {
		return common.ErrCardMethodRegister.ErrorWithExtStr("method/appID is empty")
	}
	if handler == nil {
		return common.ErrCardMethodRegister.ErrorWithExtStr("action handler is nil")
	}

	cardHandler.Set(appID, method, handler)

	return nil
}

func CardCallBack(ctx context.Context, appID string, header map[string]string, body []byte) (*protocol.CardForm, string, error) {
	// check params
	if appID == "" || len(header) == 0 || len(body) == 0 {
		return nil, "", common.ErrCardParams.ErrorWithExtStr("callBack params is empty")
	}

	// check init
	if cardHandler == nil {
		return nil, "", common.ErrCardManagerNotInit.Error()
	}

	// get app config
	appConf, err := appconfig.GetConfig(appID)
	if err != nil {
		return nil, "", common.ErrAppConfNotFound.ErrorWithExtErr(err)
	}

	// challenge event
	f := &protocol.CardChallenge{}
	err = json.Unmarshal(body, f)
	if err != nil {
		return nil, "", common.ErrJsonUnmarshal.ErrorWithExtErr(fmt.Errorf("challenge error[%v]", err))
	}

	// is challenge callback
	if len(f.Challenge) != 0 {
		if f.Token != appConf.VerifyToken {
			return nil, "", common.ErrCardVeriTokenInvalid.Error()
		}

		return nil, f.Challenge, nil
	}

	// action callback
	callback := &protocol.CardCallbackForm{}
	err = json.Unmarshal(body, callback)
	if err != nil {
		return nil, "", common.ErrJsonUnmarshal.ErrorWithExtErr(fmt.Errorf("card callback error[%v]", err))
	}

	// check signature
	if !cardHandler.ignoreSign[appID] {
		err = verifySignature(ctx, appConf.VerifyToken, header, body)
		if err != nil {
			return nil, "", common.ErrCardSignatureInvalid.Error()
		}
	}

	var ok bool
	var method string
	if method, ok = callback.Action.Value["method"]; !ok {
		return nil, "", common.ErrCardWithoutMethod.ErrorWithExtErr(err)
	}
	if _, ok := callback.Action.Value["sid"]; !ok {
		return nil, "", common.ErrCardWithoutSessionID.ErrorWithExtErr(err)
	}

	// check meta
	meta := &protocol.Meta{}
	if metaData, ok := callback.Action.Value["meta"]; !ok {
		meta = nil
	} else {
		err = json.Unmarshal([]byte(metaData), meta)
		if err != nil {
			return nil, "", common.ErrCardMetaInvalid.ErrorWithExtErr(fmt.Errorf("meta jsonUnmarshalError[%v]", err))
		}
	}

	handler, err := cardHandler.Get(appID, method)
	if err != nil {
		return nil, "", common.ErrCardMethodRegister.ErrorWithExtErr(err)
	}
	if handler == nil {
		return nil, "", common.ErrCardHandlerIsNil.ErrorWithExtStr(fmt.Sprintf("method[%s]", method))
	}

	card, err := handler(ctx, callback)
	if err != nil {
		return nil, "", common.ErrCardHandlerFailed.ErrorWithExtErr(err)
	}

	return card, "", nil
}

// check action Signature
func verifySignature(ctx context.Context, verifyToken string, header map[string]string, body []byte) error {
	timestamp := header["X-Lark-Request-Timestamp"]
	nonce := header["X-Lark-Request-Nonce"]
	sig := header["X-Lark-Signature"]

	targetSig := genPostRequestSignature(nonce, timestamp, string(body), verifyToken)
	if sig == targetSig {
		return nil
	}

	return fmt.Errorf("signature invalid")
}

func genPostRequestSignature(nonce string, timestamp string, body string, token string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(token)
	b.WriteString(body)

	bs := []byte(b.String())
	h := sha1.New()
	h.Write(bs)
	bs = h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
