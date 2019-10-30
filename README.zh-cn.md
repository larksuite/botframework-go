# botframework-go
飞书开放平台应用开发接口 golang 版本 SDK，支持开发者快速搭建和开发飞书应用。  
支持开发基于飞书开放平台的机器人应用、小程序应用。  
支持自动生成 gin 框架代码。  

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
- 通讯录变更
- 接收消息
- 机器人被移出群聊
- 机器人被邀请进入群聊
- 用户和机器人的会话首次被创建
- app_ticket事件
- 应用启动、停用事件
- 用户进群出群事件通知
- 解散群通知
- 群配置修改事件
- 应用商店应用购买
- 创建小组件实例事件
- 删除小组件实例事件

## 授权
获取 tenant_access_token (支持 ISV apps or internal apps)
- GetTenantAccessToken

获取 app_access_token (支持 ISV apps or internal apps)
- GetAppAccessToken

触发推送 app_ticket
- ReSendAppTicket

## 身份验证
- 小程序身份验证
- Open SSO 身份验证

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

## 消息构造
- 富文本构造
- 卡片构造

# 目录说明
- SDK:                封装开放平台相关通用操作
    - appconfig:      应用相关配置信息
    - auth:           封装开放平台授权相关接口
    - authentication: 封装身份认证相关接口
    - chat:           封装开放平台机器人群信息和群管理相关接口
    - common:         SDK公共操作集合
    - event:          封装事件订阅、卡片action回调、机器人接收消息回调的接口
    - message:        封装机器人发送消息的接口，支持发送文本、图片、富文本、群名片、卡片消息，支持批量发送消息，提供简单的构造富文本、卡片消息的接口。
    - protocol:       开放平台相关协议、SDK自定义协议
- generatecode:       框架代码生成工具，当前只支持生成gin框架的代码

# SDK 使用说明
## 日志初始化
SDK日志初始化函数 `func InitLogger(log LogInterface, option interface{})`  

SDK提供默认的日志实现和默认的日志参数，通过如下调用使用默认日志实现: `common.InitLogger(common.NewCommonLogger(), common.DefaultOption())`  

开发者可以通过实现`LogInterface`接口，来使用自定义日志库。  

## SDK初始化
SDK使用前需要执行初始化操作，具体操作步骤：  
1. 获取应用相关配置信息（AppID、AppSecret、VerifyToken、EncryptKey 等），为了数据安全，不建议开发者在代码中明文写入这些信息，您可以选择从数据库读取、远程配置系统、环境变量中获取；  
2. 调用`appconfig.Init(conf)`函数，初始化应用配置；  
3. 如果是 Independent Software Vendor App (ISVApp), 需要实现读取和保存 AppTicket 的接口。框架提供使用redis读取和保存 AppTicket 的接口实现，你也可以实现`TicketManager`接口，来使自定义 AppTicket 的读写方式;  
4. 根据业务逻辑注册事件回调处理函数;  

### 示例代码
1. 从redis读取配置信息：  
   [示例代码](./demo/sdk_init/sdk_init.go)
2. 从环境变量读取配置信息：  
```golang
conf := &appconfig.AppConfig{
    AppID: os.Getenv("AppID"),//从飞书开放平台-凭证与基础信息中获取
    AppType: protocol.InternalApp, //apptype只有两种，Independent Software Vendor App 和 Internal App
    AppSecret:   os.Getenv("AppSecret"),//从飞书开放平台的凭证与基础信息中获取
    VerifyToken: os.Getenv("VerifyToken"),//从飞书开放平台的事件订阅中获取
    EncryptKey:  os.Getenv("EncryptKey"),//从飞书开放平台的事件订阅中获取
}
```  

## 存储初始化
SDK存储读写接口`type DBClient interface`。在SDK中，其只要用于： ISV 应用的 app ticket 读写操作；身份认证操作的session数据读写操作；敏感信息如 app secret的读取。  

接口提供默认的redis实现`DefaultRedisClient`，使用前需要做初始化操作，初始化示例代码  
```golang
redisClient := &common.DefaultRedisClient{}
err := redisClient.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
if err != nil {
    return fmt.Errorf("init redis-client error[%v]", err)
}
```  

开发者可以通过实现 DBClient 接口，来使用自定义存储。  

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
- 首次生成代码时，会依据配置文件生成全部代码文件；  
- 之后若修改配置文件（修改代码路径之外的其他选项），在原始的路径上重新生成代码时，只会强制更新`./handler/regist.go`文件，其他文件不会更新，以避免覆盖用户自定义代码。  
- `./handler/regist.go`文件，会被强制更新，用户不应该在该文件中加入自定义代码。  

