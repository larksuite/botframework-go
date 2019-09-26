# botframework-go
botframework-go is a golang Software Development Kit which gives tools and interfaces to developers needed to access to lark open platform APIs.
It support to generate code using Gin framework.

# List of Main Features
## Subscribe to Event Notification
Interfaces
- EventRegister
- EventCallback

List of supported notification
- approval
- leave approval
- work approval
- shift approval
- remedy approval
- trip approval
- app open
- contacts change
- message
- remove bot
- add bot
- create p2p chat
- app ticket
- app status change

## Authorization
Obtain tenant_access_token (ISV apps or internal apps)
- GetTenantAccessToken

Obtain app_access_token (ISV apps or internal apps)
- GetAppAccessToken

Re-pushing app_ticket
- ReSendAppTicket

## Bot Send Message
Interfaces
- SendTextMessage
- SendImageMessage
- SendRichTextMessage
- SendShareChatMessage
- SendCardMessage
- SendTextMessageBatch
- SendImageMessageBatch
- SendRichTextMessageBatch
- SendShareChatMessageBatch
- SendCardMessageBatch

## Group
Interfaces
- GetChatInfo
- GetChatList
- CheckUserInGroup
- CheckBotInGroup
- CheckUserBotInSameGroup

## Richtext Builder
build richtext demo
```go
postForm := make(map[protocol.Language]*protocol.RichTextForm)

// en-us
titleUS := "this is a title"
contentUS := message.NewRichTextContent()

// first line
contentUS.AddElementBlock(
    message.NewTextTag("first line: ", true, 1),
    message.NewATag("hyperlinks ", true, "https://www.feishu.cn"),
    message.NewAtTag("username", userID),
)

// second line
contentUS.AddElementBlock(
    message.NewTextTag("second line: ", true, 1),
    message.NewTextTag("text test", true, 1),
)

postForm[protocol.EnUS] = message.NewRichTextForm(&titleUS, contentUS)

// zh-cn
titleCN := "这是一个标题"
contentCN := message.NewRichTextContent()

// first line
contentCN.AddElementBlock(
    message.NewTextTag("第一行: ", true, 1),
    message.NewATag("超链接 ", true, "https://www.feishu.cn"),
    message.NewAtTag("username", userID),
)

// second line
contentCN.AddElementBlock(
    message.NewTextTag("第二行: ", true, 1),
    message.NewTextTag("文本测试", true, 1),
)

postForm[protocol.ZhCN] = message.NewRichTextForm(&titleCN, contentCN)
```
See more examples in file "SDK/message/message_test.go"

## Card Builder
build Card demo
```go
//card builder
builder := &message.CardBuilder{}

//add config
config := protocol.ConfigForm{
    MinVersion:     protocol.VersionForm{},
    WideScreenMode: true,
}
builder.SetConfig(config)

//add header
content := "Please choose color"
line := 1
title := protocol.TextForm{
    Tag:     protocol.PLAIN_TEXT_E,
    Content: &content,
    Lines:   &line,
}
builder.AddHeader(title, "")

//add hr
builder.AddHRBlock()

//add block
builder.AddDIVBlock(nil, []protocol.FieldForm{
    *message.NewField(false, message.NewMDText("**Async**", nil, nil, nil)),
}, nil)

//add divBlock
builder.AddDIVBlock(nil, []protocol.FieldForm{
    *message.NewField(false, message.NewMDText("**Sync**", nil, nil, nil)),
}, nil)

//add actionBlock
payload1 := make(map[string]string, 0)
payload1["color"] = "red"
builder.AddActionBlock([]protocol.ActionElement{
    message.NewButton(message.NewMDText("red", nil, nil, nil),
        nil, nil, payload1, protocol.PRIMARY, nil, "asyncButton"),
})

//add jumpBlock
url := "https://www.google.com"
ext := message.NewJumpButton(message.NewMDText("jump to google", nil, nil, nil), &url, nil, protocol.DEFAULT)
builder.AddDIVBlock(message.NewMDText("", nil, nil, nil), nil, ext)

//add imageBlock
builder.AddImageBlock(
    message.NewMDText("", nil, nil, nil),
    *message.NewMDText("", nil, nil, nil),
    imageKey)

//generate card
card, err := builder.BuildForm()
```
See more examples in file "SDK/message/message_test.go"

# Directory Description
- SDK           : Lark open platform APIs
    - appconfig : Appinfo config
    - auth      : Authorization
    - chat      : Group
    - common    : Common functions/definition
    - event     : Event notification/card action callback/bot command callback
    - message   : Bot send message。
    - protocol  : Lark open platform protocol
- generatecode  : Generate code using Gin framework

# Generate code using Gin framework
## Config
```yml
ServiceInfo:
  Path: github.com/larksuite/demo  # relative path in GOPATH or go module name
  GenCodePath:  # Code generation absolute path. If it is empty, the configuration item named "Path" will be used.
  EventWebhook: /webhook/event
  CardWebhook: /webhook/card
  AppID: cli_12345            # your app id
  Description: test_demo      # your app description
  IsISVApp: false             # ISV App flag, false is default
EventList:
  - EventName: Message           # required
  # - EventName: AppTicket       # use as needed, ISVApp must
  # - EventName: Approval        # use as needed
  # - EventName: LeaveApproval   # use as needed
  # - EventName: WorkApproval    # use as needed
  # - EventName: ShiftApproval   # use as needed
  # - EventName: RemedyApproval  # use as needed
  # - EventName: TripApproval    # use as needed
  # - EventName: AppOpen         # use as needed
  # - EventName: ContactUser     # use as needed
  # - EventName: ContactDept     # use as needed
  # - EventName: ContactScope    # use as needed
  # - EventName: RemoveBot       # use as needed
  # - EventName: AddBot          # use as needed
  # - EventName: P2PChatCreate   # use as needed
  # - EventName: AppStatusChange # use as needed
CommandList:
  - Cmd: default # required
    Description: Text that is empty or isnot matched
#   - Cmd: help
#     Description: Text that begin with the word help
#   - Cmd: show
#     Description: Text that begin with the word show
CardActionList:
  - MethodName: create
  - MethodName: delete
  - MethodName: update
```

## Command
```shell
# cd projectPath
go build
./botframework-go -f ./generatecode/demo.yml
```

## Generate Code Rule
If the code is first generated, all code files are generated by the configuration file.
If you modify the configuration file later, and regenerate the code on the original path,
only the "./handler/regist.go" file will be forced updated, other files are not updated to avoid overwriting user-defined code.
Because the "./handler/regist.go" file will be forced update, you should not write your business code in the file.
