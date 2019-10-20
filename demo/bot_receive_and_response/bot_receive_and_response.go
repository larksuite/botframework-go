package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/event"
	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
	"github.com/larksuite/botframework-go/demo/sdk_init"
)

const (
	appID     = ""
	tenantKey = ""
)

func main() {
	r := gin.Default()

	common.InitLogger(common.NewCommonLogger(), common.DefaultOption())
	defer common.FlushLogger()

	//Necessary steps，init app configuration and regist
	err := InitInfoAndRegist()
	if err != nil {
		common.Logger(context.TODO()).Errorf("InitError[%v]", err)
		return
	}

	r.POST("/webhook/event", EventCallback) //open platform event callback
	r.POST("/webhook/card", CardCallback)   //card action callback

	r.Run(":8089")
}

//init app configuration and regist
func InitInfoAndRegist() error {
	err := sdk_init.InitInfo()
	if err != nil {
		common.Logger(context.TODO()).Errorf("InitError[%v]", err)
		return err
	}
	RegistHandler(appID)
	return nil
}

func RegistHandler(appID string) {
	event.EventRegister(appID, protocol.EventTypeMessage, EventMessage) //necessary function.process events and distribute them
	event.BotRecvMsgRegister(appID, "help", BotRecvMsgHelp)             //response "help"
	event.BotRecvMsgRegister(appID, "card", BotRecvMsgCard)             //response "card"
	event.CardRegister(appID, "clickbutton", ActionClickButton)         //response when clicking this button
}

// EventCallback open platform event
func EventCallback(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		common.Logger(c).Errorf("eventReqParamsError: readHttpBodyError err[%v]bodyLen[%d]", err, len(body))
		c.JSON(500, gin.H{"codemsg": common.ErrEventParams.String()})
		return
	}

	appID := appID
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

	appID := appID
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

// event-Message
func EventMessage(ctx context.Context, eventBody []byte) error {
	return event.BotRecvMsgHandler(ctx, eventBody)
}

func BotRecvMsgHelp(ctx context.Context, msg *protocol.BotRecvMsg) error {
	user := &protocol.UserInfo{
		ID:   msg.OpenChatID,
		Type: protocol.UserTypeChatID,
	}
	message.SendTextMessage(ctx, tenantKey, appID, user, "", "hello,this is help")
	return nil
}

func BotRecvMsgCard(ctx context.Context, msg *protocol.BotRecvMsg) error {
	user := &protocol.UserInfo{
		ID:   msg.OpenChatID,
		Type: protocol.UserTypeChatID,
	}

	//build card
	builder := &message.CardBuilder{}
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)

	//add header
	content := "this is a card"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")

	//add button
	button1 := make(map[string]string, 0)
	button1["key"] = "buttonValue"
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewButton(message.NewMDText("callback button", nil, nil, nil), nil, nil, button1, protocol.DANGER, nil, "clickbutton"),
	})

	card, err := builder.BuildForm()
	if err != nil {
		common.Logger(ctx).Errorf("card build failed error[%v]", err)
		return fmt.Errorf("card build failed error[%v]", err)
	}

	//add params to use message.SendCardMessage
	resp, err := message.SendCardMessage(ctx, tenantKey, appID, user, "", *card, false)
	if err != nil {
		common.Logger(ctx).Errorf("send message failed error[%v]", err)
		return fmt.Errorf("send message failed error[%v]", err)
	}

	common.Logger(ctx).Errorf("code[%d],msg[%s],openMessageID[%s]", resp.Code, resp.Msg, resp.Data.MessageID)
	return nil
}

// methodName-clickbutton
func ActionClickButton(ctx context.Context, callback *protocol.CardCallbackForm) (*protocol.CardForm, error) {
	method, _ := callback.Action.Value["method"]
	sessionID, _ := callback.Action.Value["sid"]
	common.Logger(ctx).Infof("cardActionCallBack: method[%s]sessionID[%s]", method, sessionID)

	//update card
	//build a new card
	builder := &message.CardBuilder{}
	//add config
	config := protocol.ConfigForm{
		MinVersion:     protocol.VersionForm{},
		WideScreenMode: true,
	}
	builder.SetConfig(config)

	//add header
	content := "this is a card"
	line := 1
	title := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &content,
		Lines:   &line,
	}
	builder.AddHeader(title, "")

	//add button
	button1 := make(map[string]string, 0)
	button1["key"] = "buttonValue"
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewButton(message.NewMDText("clicked button", nil, nil, nil), nil, nil, button1, protocol.UNKNOWN, nil, "clickbutton"),
	})

	card, err := builder.BuildForm()
	if err != nil {
		common.Logger(ctx).Errorf("card update failed error[%v]", err)
		return &protocol.CardForm{}, fmt.Errorf("card update failed error[%v]", err)
	}

	return card, nil
}
