# botframework-go
**Deprecated, please use [oapi-sdk-go](https://github.com/larksuite/oapi-sdk-go)**

botframework-go is a golang Software Development Kit which gives tools and interfaces to developers needed to access to Lark Open Platfom APIs.  
It support for developing a robot application or mini-program application based on the Lark Open Platfom.  
It support for generateing code using Gin framework.  

# List of Main Features
## Subscribe to Event Notification
Interfaces
- EventRegister
- EventCallback

List of supported notification
- Request approved
- Leave approved
- Overtime approved
- Shift change approved
- Correction request approved
- Business trip approved
- App enabled
- Contacts updates
- Message
- Bot removed from group chat
- Bot invited to group chat
- Create p2p chat
- App ticket event
- App status change
- User in or out the group chat
- Disband chat
- Change group chat config
- Order paid
- Create a widget instance event
- Delete a widget instance event
- Read message
- App uninstall

## Authorization
Obtain tenant_access_token (ISV apps or internal apps)
- GetTenantAccessToken

Obtain app_access_token (ISV apps or internal apps)
- GetAppAccessToken

Re-pushing app_ticket
- ReSendAppTicket

## Authentication
- Mini Program Authentication
- Open SSO Authentication

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
- UpdateChatInfo
- CreateChat
- AddUserToChat
- DeleteUserFromChat
- DisbandChat

## Message Builder
- Richtext Builder
- Card Builder

# Directory Description
- SDK:                Lark open platform APIs
    - appconfig:      Appinfo config
    - auth:           Authorization
    - authentication: Authentication
    - chat:           Group
    - common:         Common functions/definition
    - event:          Event notification/card action callback/bot command callback
    - message:        Bot send message
    - protocol:       Lark open platform protocol
- generatecode:       Generate code using Gin framework

# SDK Instruction
## Log Initialization
Log Initialization function `func InitLogger(log LogInterface, option interface{})`. demo: `common.InitLogger(common.NewCommonLogger(), common.DefaultOption())`  

Developers can use the custom log library by implementing the LogInterface interface  

## SDK Initialization
You must init app config before using SDK. You can follow these steps to do it.  
1. Get necessary information, such as AppID, AppSecret, VerifyToken and EncryptKey. But you'd better not use clear text in code. Getting them from database or environment variable is a better way.  
2. Init app config.  
3. Get appTickey if your app is Independent Software Vendor App, otherwise you can ignore this.  
4. Register events what you want.  

### Example code  
1. get config from redis.  
   [demo code](./demo/sdk_init/sdk_init.go)
2. get config from environment variable.  
```go
conf := &appconfig.AppConfig{
    AppID: os.Getenv("AppID"),              //get it from lark-voucher and basic information。
    AppType: protocol.InternalApp,          //AppType only has two types: Independent Software Vendor App（ISVApp） or Internal App.
    AppSecret:   os.Getenv("AppSecret"),    //get it from lark-voucher and basic information.
    VerifyToken: os.Getenv("VerifyToken"),  //get it from lark-event subscriptions.
    EncryptKey:  os.Getenv("EncryptKey"),   //get it from lark-event subscriptions.
}
```

## DB-client Initialization
`DBClient` SDK DB Client interface.  
It is used for: read/write app ticket ; read/write user session data for authentication operations; read sensitive information such as app secret.  

`DefaultRedisClient` Default Redis Client. Need to be initialized before use.  
demo:  
```golang
redisClient := &common.DefaultRedisClient{}
err := redisClient.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
if err != nil {
    return fmt.Errorf("init redis-client error[%v]", err)
}
```  

Developers can use the custom db client library by implementing the DBClient interface  

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
  - EventName: Message                # required
  # - EventName: AppTicket            # use as needed, ISVApp must
  # - EventName: Approval             # use as needed
  # - EventName: LeaveApproval        # use as needed
  # - EventName: WorkApproval         # use as needed
  # - EventName: ShiftApproval        # use as needed
  # - EventName: RemedyApproval       # use as needed
  # - EventName: TripApproval         # use as needed
  # - EventName: AppOpen              # use as needed
  # - EventName: ContactUser          # use as needed
  # - EventName: ContactDept          # use as needed
  # - EventName: ContactScope         # use as needed
  # - EventName: RemoveBot            # use as needed
  # - EventName: AddBot               # use as needed
  # - EventName: P2PChatCreate        # use as needed
  # - EventName: AppStatusChange      # use as needed
  # - EventName: UserToChat           # use as needed
  # - EventName: ChatDisband          # use as needed
  # - EventName: GroupSettingUpdate   # use as needed
  # - EventName: OrderPaid            # use as needed
  # - EventName: CreateWidgetInstance # use as needed
  # - EventName: DeleteWidgetInstance # use as needed
  # - EventName: MessageRead          # use as needed
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
- If the code is first generated, all code files are generated by the configuration file.  
- If you modify the configuration file later, and regenerate the code on the original path, only the `./handler/regist.go` file will be forced updated, other files are not updated to avoid overwriting user-defined code.  
- Because the `./handler/regist.go` file will be forced update, you should not write your business code in the file.  

# demo
- [Manage AppAccessToken](./docs/us/app_access_token.md)
- [Manage TenantAccessToken](./docs/us/tenant_access_token.md)
- [Webhook_Event](./docs/us/webhook_event.md)
- [Webhook_Card](./docs/us/webhook_card.md)
- [Send Message](./docs/us/send_message.md)
- [Build Card](./docs/us/send_card.md)
- [Authentication](./docs/us/authentication.md)
