# botframework-go
飞书开放平台应用开发接口golang版本SDK，支持开发者快速搭建和开发飞书应用,支持自动生成gin框架代码。

# 支持接口列表
## 订阅事件通知
接口：
- EventRegister
- EventCallback

支持事件列表:
- 审批通过
- 请假审批
- 加班审批
- 换班审批
- 补卡审批
- 出差审批
- 开通应用
- 通讯录变更-用户变更
- 通讯录变更-部门变更
- 通讯录变更-变更权限范围
- 接收消息
- 机器人被移出群聊
- 机器人被邀请进入群聊
- 用户和机器人的会话首次被创建
- app_ticket事件
- 应用启动、停用事件

## 授权
获取 tenant_access_token (支持 ISV apps or internal apps)
- GetTenantAccessToken

获取 app_access_token (支持 ISV apps or internal apps)
- GetAppAccessToken

触发推送 app_ticket
- ReSendAppTicket

## 机器人发送消息
接口
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


## 群信息和群管理
接口
- GetChatInfo
- GetChatList
- CheckUserInGroup
- CheckBotInGroup
- CheckUserBotInSameGroup

## 富文本构造
构造富文本代码示例
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
之后, 可以调用SendRichTextMessage函数发送富文本信息

## 卡片构造
构造卡片代码事例
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
更多的卡片构造请查看:SDK/message/card_builder_test.go

# 目录说明
- SDK           : 封装开放平台相关通用操作
    - appconfig : 应用相关配置信息
    - auth      : 封装开放平台授权相关接口
    - chat      : 封装开放平台机器人群信息和群管理相关接口
    - common    : SDK公共操作集合
    - event     : 封装事件订阅、卡片action回调、机器人接收消息回调的接口
    - message   : 封装机器人发送消息的接口，支持发送文本、图片、富文本、群名片、卡片消息，支持批量发送消息，提供简单的构造富文本、卡片消息的接口。
    - protocol  : 开放平台相关协议、SDK自定义协议
- generatecode: 框架代码生成工具，当前只支持生成gin框架的代码

# 生成 Gin 框架代码
## 配置文件示例
```yml
ServiceInfo:
  Path: github.com/larksuite/demo  # GOPATH相对路径，或者使用go module方式时的module name
  GenCodePath:  # 生成代码的绝对路径；若为空，代码会生成到配置项Path对应的GOPATH路径下
  EventWebhook: /webhook/event
  CardWebhook: /webhook/card
  AppID: cli_12345        # 应用ID
  Description: test_demo  # 应用描述信息
  IsISVApp: false         # ISV 应用标志, 默认为非ISV应用
EventList:
  - EventName: Message           # 必须
  # - EventName: AppTicket       # 按需使用, ISV应用 必须订阅
  # - EventName: Approval        # 按需使用
  # - EventName: LeaveApproval   # 按需使用
  # - EventName: WorkApproval    # 按需使用
  # - EventName: ShiftApproval   # 按需使用
  # - EventName: RemedyApproval  # 按需使用
  # - EventName: TripApproval    # 按需使用
  # - EventName: AppOpen         # 按需使用
  # - EventName: ContactUser     # 按需使用
  # - EventName: ContactDept     # 按需使用
  # - EventName: ContactScope    # 按需使用
  # - EventName: RemoveBot       # 按需使用
  # - EventName: AddBot          # 按需使用
  # - EventName: P2PChatCreate   # 按需使用
  # - EventName: AppStatusChange # 按需使用
CommandList:
  - Cmd: Default # 必须
    Description: 表示默认命令，群聊只@机器人而不输入任何其他内容，或收到未定义的命令时
  # - Cmd: Help
  #   Description: 向机器人发送消息，前缀带有help
  # - Cmd: Show
  #   Description: 向机器人发送消息，前缀带有show
CardActionList:
  - MethodName: create
  - MethodName: delete
  - MethodName: update
```

## 生成代码命令
```shell
# cd projectPath
go build
./botframework-go -f ./generatecode/demo.yml
```

## 生成代码规则说明
- 首次生成代码时，会依据配置文件生成全部代码文件；之后若修改配置文件（修改代码路径之外的其他选项），在原始的路径上重新生成代码时，只会强制更新`./handler/regist.go`文件，其他文件不会更新，以避免覆盖用户自定义代码。`./handler/regist.go`文件，会被强制更新，用户不应该在该文件中加入自定义代码。
