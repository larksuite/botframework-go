// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import "fmt"

type ErrCodeMsg struct {
	Code    int
	Message string
}

func (e ErrCodeMsg) Error() error {
	return fmt.Errorf("code=%d, msg=%s", e.Code, e.Message)
}

func (e ErrCodeMsg) ErrorWithExtErr(err error) error {
	return fmt.Errorf("code=%d, msg=%s, extError=%v", e.Code, e.Message, err)
}

func (e ErrCodeMsg) ErrorWithExtStr(str string) error {
	return fmt.Errorf("code=%d, msg=%s, extError=%s", e.Code, e.Message, str)
}

func (e ErrCodeMsg) String() string {
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Message)
}

func (e ErrCodeMsg) StringWithExtErr(extErr error) string {
	return fmt.Sprintf("code=%d, msg=%s, extError=%v", e.Code, e.Message, extErr)
}

var (
	Success = &ErrCodeMsg{Code: 0, Message: "success"}

	// 1. common 1000 - 1999
	ErrJsonMarshal        = &ErrCodeMsg{Code: 1000, Message: "json marshal error"}
	ErrJsonUnmarshal      = &ErrCodeMsg{Code: 1001, Message: "json unmarshal error"}
	ErrOpenApiFailed      = &ErrCodeMsg{Code: 1002, Message: "open api failed"}
	ErrOpenApiReturnError = &ErrCodeMsg{Code: 1003, Message: "open api return error"}
	ErrAppConfNotFound    = &ErrCodeMsg{Code: 1004, Message: "app config not found"}

	// 2. auth 2000 - 2999
	ErrAppTokenNotFound             = &ErrCodeMsg{Code: 2000, Message: "auth app token not found"}
	ErrAppTicketNotFound            = &ErrCodeMsg{Code: 2001, Message: "auth app ticket not found"}
	ErrGetInternalTenantAccessToken = &ErrCodeMsg{Code: 2002, Message: "auth get internal tenant access token error"}
	ErrGetISVTenantAccessToken      = &ErrCodeMsg{Code: 2003, Message: "auth get ISV tenant access token error"}
	ErrGetAppAccessToken            = &ErrCodeMsg{Code: 2004, Message: "auth get app access token error"}
	ErrGetInternalAppAccessToken    = &ErrCodeMsg{Code: 2005, Message: "auth get internal app access token error"}
	ErrGetISVAppAccessToken         = &ErrCodeMsg{Code: 2006, Message: "auth get ISV app access token error"}
	ErrRespDataIsNil                = &ErrCodeMsg{Code: 2007, Message: "auth response data is nil"}
	ErrTicketManagerNotInit         = &ErrCodeMsg{Code: 2008, Message: "auth ticket manager not init"}
	ErrSetAppTicketFailed           = &ErrCodeMsg{Code: 2009, Message: "auth set app ticket failed"}

	// 3. message 3000 - 3999
	ErrSendMsgParams     = &ErrCodeMsg{Code: 3000, Message: "send msg params error"}
	ErrPostFormParams    = &ErrCodeMsg{Code: 3001, Message: "postform params error"}
	ErrImageParams       = &ErrCodeMsg{Code: 3001, Message: "postform params error"}
	ErrGenBinImageFailed = &ErrCodeMsg{Code: 3001, Message: "generate binary image error"}

	ErrCardUpdateParams = &ErrCodeMsg{Code: 3100, Message: "update card params error"}

	// 4. chat 4000 - 4999
	ErrChatParams = &ErrCodeMsg{Code: 4000, Message: "chat params error"}

	// 5. event 5000 - 5999
	ErrEventTypeRegister      = &ErrCodeMsg{Code: 5000, Message: "event type register handler error"}
	ErrEventManagerNotInit    = &ErrCodeMsg{Code: 5001, Message: "event not init"}
	ErrEventParams            = &ErrCodeMsg{Code: 5002, Message: "event params error"}
	ErrEventDecrypt           = &ErrCodeMsg{Code: 5003, Message: "event decrypt error"}
	ErrEventGetBase           = &ErrCodeMsg{Code: 5004, Message: "event get event base error"}
	ErrEventVeriToken         = &ErrCodeMsg{Code: 5005, Message: "event veri token error"}
	ErrEventTypeUnknown       = &ErrCodeMsg{Code: 5006, Message: "event unknown callback type"}
	ErrEventGetJsonEvent      = &ErrCodeMsg{Code: 5007, Message: "event get event body error"}
	ErrEventGetJsonType       = &ErrCodeMsg{Code: 5008, Message: "event get event type error"}
	ErrEventGetJsonAppID      = &ErrCodeMsg{Code: 5009, Message: "event get app id error"}
	ErrEventAppIDNotMatch     = &ErrCodeMsg{Code: 5010, Message: "event app id not match error"}
	ErrEventAppIDUnregistered = &ErrCodeMsg{Code: 5011, Message: "event appid_handler unregistered"}
	ErrEventTypeUnregistered  = &ErrCodeMsg{Code: 5012, Message: "event notification_type_handler unregistered"}
	ErrEventHandlerIsNil      = &ErrCodeMsg{Code: 5013, Message: "event handler not found"}
	ErrEventHandlerFailed     = &ErrCodeMsg{Code: 5014, Message: "event handle function error"}

	ErrBotRecvMsgRegister       = &ErrCodeMsg{Code: 5100, Message: "botRecvMsg registered error"}
	ErrBotRecvMsgMsgTypeJson    = &ErrCodeMsg{Code: 5101, Message: "botRecvMsg get msg_type error"}
	ErrBotRecvMsgAppIDJson      = &ErrCodeMsg{Code: 5102, Message: "botRecvMsg get app_id error"}
	ErrBotRecvMsgHandlerNoFound = &ErrCodeMsg{Code: 5103, Message: "botRecvMsg cannot find handler"}
	ErrBotRecvMsgHandlerFailed  = &ErrCodeMsg{Code: 5104, Message: "botRecvMsg call handler failed"}

	ErrCardParams           = &ErrCodeMsg{Code: 5200, Message: "card action callback params error"}
	ErrCardMethodRegister   = &ErrCodeMsg{Code: 5201, Message: "card action method has not registered yet"}
	ErrCardManagerNotInit   = &ErrCodeMsg{Code: 5202, Message: "card action handler need be init"}
	ErrCardVeriTokenInvalid = &ErrCodeMsg{Code: 5203, Message: "card action callback veri token invalid"}
	ErrCardSignatureInvalid = &ErrCodeMsg{Code: 5204, Message: "card action callback signature invalid"}
	ErrCardWithoutMethod    = &ErrCodeMsg{Code: 5205, Message: "card action callback has no method"}
	ErrCardWithoutSessionID = &ErrCodeMsg{Code: 5206, Message: "card action callback has no session id"}
	ErrCardMetaInvalid      = &ErrCodeMsg{Code: 5207, Message: "card action callback meta invalid"}
	ErrCardHandlerIsNil     = &ErrCodeMsg{Code: 5208, Message: "card action handler not found"}
	ErrCardHandlerFailed    = &ErrCodeMsg{Code: 5209, Message: "card action handler failed"}

	// 6. authentication 6000 - 6999
	ErrValidateParams     = &ErrCodeMsg{Code: 6000, Message: "authentication-login params error"}
	ErrAuthParams         = &ErrCodeMsg{Code: 6001, Message: "authentication-auth  params error"}
	ErrMinaCodeGetParams  = &ErrCodeMsg{Code: 6002, Message: "mini-program codeToSession get params error"}
	ErrMinaSetAuth        = &ErrCodeMsg{Code: 6003, Message: "mini-program set auth-user-info error"}
	ErrMinaGetAuth        = &ErrCodeMsg{Code: 6004, Message: "mini-program get auth-user-info error"}
	ErrMinaSessionInvalid = &ErrCodeMsg{Code: 6005, Message: "mini-program session invalid"}
	ErrLogoutParams       = &ErrCodeMsg{Code: 6006, Message: "authentication-logout  params error"}
)
