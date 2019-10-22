package main

import (
	"context"
	"fmt"
	"time"

	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
	"github.com/larksuite/botframework-go/demo/sdk_init"
)

func main() {
	//init log
	common.InitLogger(common.NewCommonLogger(), common.DefaultOption())
	defer common.FlushLogger()

	//param
	ctx := context.TODO()

	//Necessary step: init app configuration
	err := sdk_init.InitInfo()
	if err != nil {
		common.Logger(ctx).Errorf("InitError[%v]", err)
		return
	}

	//params
	chatID := ""    //p2p or group chat ID
	tenantKey := "" //tenantKey of your company
	appID := ""     //APP ID

	//send card
	sendCard(chatID, tenantKey, appID)
}

func sendCard(chatID, tenantKey, appID string) error {
	user := &protocol.UserInfo{
		ID:   chatID,
		Type: protocol.UserTypeChatID,
	}

	ctx := context.TODO()

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

	//add dividing line
	builder.AddHRBlock()

	//set href
	urlhref := make(map[string]protocol.URLForm, 0)
	uHost := "https://www.larksuite.com/"
	urlhref["link"] = protocol.URLForm{Url: &uHost}

	//add content block,suck as title、content and extra（null here）.
	builder.AddDIVBlock(message.NewMDText("[lark]($link)", nil, nil, urlhref), []protocol.FieldForm{
		*message.NewField(false, message.NewMDText("**boldText**", nil, nil, nil)),
		*message.NewField(false, message.NewMDText("text", nil, nil, nil)),
	}, nil)

	builder.AddHRBlock()

	//add image block
	ImageContent := "Description when your mouse is over the picture"
	ImageTag := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &ImageContent,
		Lines:   nil,
	}

	//get imageKey.You can get detailed introduction about imageKey from file(6.2)
	imageURL := "https://is3-ssl.mzstatic.com/image/thumb/Purple113/v4/ed/ee/c0/edeec03e-d111-ac8d-3441-409acd11dbea/source/512x512bb.jpg"
	imageKey, err := message.GetImageKey(ctx, tenantKey, appID, imageURL, "")
	if err != nil {
		common.Logger(ctx).Errorf("get imageKey failed[%v]", err)
	}
	builder.AddImageBlock(message.NewMDText("image title", nil, nil, nil), ImageTag, imageKey)

	builder.AddHRBlock()

	//add button
	payload1 := make(map[string]string, 0)
	payload1["key"] = "value"
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewButton(message.NewMDText("lark button", nil, nil, nil), &uHost, nil, payload1, protocol.PRIMARY, nil, "testcard"),
		message.NewButton(message.NewMDText("simple button", nil, nil, nil), nil, nil, payload1, protocol.DANGER, nil, "testcard"),
	})
	builder.AddHRBlock()

	//add menu
	//SelectPersonMenu
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewSelectPersonMenu(message.NewMDText("SelectPersonMenu", nil, nil, nil), nil, []protocol.OptionForm{}, nil, nil, "testcard"),
	},
	)

	//SelectStaticMenu
	optiontextcontent1 := "option1"
	optiontext1 := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &optiontextcontent1,
		Lines:   nil,
	}
	optiontextcontent2 := "option2"
	optiontext2 := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &optiontextcontent2,
		Lines:   nil,
	}
	value1 := "value1"
	value2 := "value2"
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewSelectStaticMenu(message.NewMDText("SelectStaticMenu", nil, nil, nil), nil, []protocol.OptionForm{
			message.NewOption(optiontext1, value1),
			message.NewOption(optiontext2, value2),
		}, &value1, nil, "testcard"), //option1 is default option
	},
	)
	builder.AddHRBlock()

	//add overflow
	overFlowcontent1 := "option1"
	overFlowtext1 := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &overFlowcontent1,
		Lines:   nil,
	}
	overFlowcontent2 := "option2"
	overFlowtext2 := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &overFlowcontent2,
		Lines:   nil,
	}
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewOverflowMenu(nil, []protocol.OptionForm{
			message.NewOption(overFlowtext1, value1),
			message.NewOption(overFlowtext2, value2),
		}, nil, "testcard"),
	},
	)
	builder.AddHRBlock()

	//add datapicker
	timePicker := time.Now().Format("2006-01-02")
	builder.AddActionBlock([]protocol.ActionElement{
		message.NewPickerDate(message.NewMDText("dataPicker", nil, nil, nil), nil, nil, &timePicker, "testcard"),
	},
	)
	builder.AddHRBlock()

	//add note
	noteTextContent := "noteTextContent"
	notetext := protocol.TextForm{
		Tag:     protocol.PLAIN_TEXT_E,
		Content: &noteTextContent,
		Lines:   nil,
	}
	noteImageContent := "noteImageContent"
	noteImage := protocol.ImageForm{
		Tag:      protocol.IMG_BLOCK,
		ImageKey: imageKey,
		ALT: protocol.TextForm{
			Tag:     protocol.PLAIN_TEXT_E,
			Content: &noteImageContent,
			Lines:   nil,
		},
	}
	builder.AddNoteBlock([]protocol.BaseElement{
		&notetext,
		&noteImage,
	})

	card, err := builder.BuildForm()
	if err != nil {
		common.Logger(ctx).Errorf("buildForm failed[%v]", err)
		return fmt.Errorf("buildForm failed[%v]", err)
	}

	//add params to use message.SendCardMessage
	resp, err := message.SendCardMessage(ctx, tenantKey, appID, user, "", *card, false)
	if err != nil {
		common.Logger(ctx).Errorf("send card failed[%v]", err)
		return fmt.Errorf("send card failed[%v]", err)
	}

	common.Logger(ctx).Info("code:[%d],msg:[%s],openMessageID:[%s]", resp.Code, resp.Msg, resp.Data.MessageID)
	return nil
}
