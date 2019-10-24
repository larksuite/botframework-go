package sdk_init

import (
	"fmt"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

//Initialization function
func InitInfo() error {
	//init log
	common.InitLogger(common.NewCommonLogger(), common.DefaultOption())
	defer common.FlushLogger()

	//init redis-client
	redisClient := &common.DefaultRedisClient{}
	err := redisClient.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
	if err != nil {
		return fmt.Errorf("init redis-client error[%v]", err)
	}

	// get app info
	conf, err := GetAppConf(redisClient)
	if err != nil {
		return fmt.Errorf("get conf error[%v]", err)
	}
	// init app info
	appconfig.Init(*conf)

	//Independent Software Vendor App（ISVApp） has to get APPTicket,if your AppType is not ISV,you can ignore this
	// ISVApp Set TicketManager
	if conf.AppType == protocol.ISVApp {
		// ISVApp need to implement the TicketManager interface
		// It is recommended to set/get your app-ticket in redis

		err := auth.InitISVAppTicketManager(auth.NewDefaultAppTicketManager(redisClient))
		if err != nil {
			return fmt.Errorf("Authorization Initialize Error[%v]", err)
		}
	}

	//If you need to register an/a event/card callback，you can do it here.
	//You can reference this simple example or detailed introduction from file（4.webhook-event和5.webhook-card）

	// regist open platform event handler
	// event.EventRegister(appID, protocol.EventTypeMessage, EventMessage)

	// regist bot recv message handler
	//event.BotRecvMsgRegister(appID, "help", BotRecvMsgHelp)

	// regist card action handler
	//event.CardRegister(appID, "testcard", ActionTestCard)

	return nil
}

// get appinfo(app_id、app_secret、veri_token、encrypt_key) from redis/mysql or remote config system
// redis/mysql or remote config system is recommended
func GetAppConf(client common.DBClient) (*appconfig.AppConfig, error) {
	//Read information from  redis
	appID, err := client.Get("AppID")
	if err != nil {
		return nil, fmt.Errorf("get AppID failed[%v]", err)
	}
	appSecret, err := client.Get("AppSecret")
	if err != nil {
		return nil, fmt.Errorf("get AppSecret failed[%v]", err)
	}
	verifyToken, err := client.Get("VerifyToken")
	if err != nil {
		return nil, fmt.Errorf("get VerifyToken failed[%v]", err)
	}
	encryptKey, err := client.Get("EncryptKey")
	if err != nil {
		return nil, fmt.Errorf("get EncryptKey failed[%v]", err)
	}

	// Initialize app config
	// Clear text in code is not recommended.
	conf := &appconfig.AppConfig{
		AppID:       appID,                //get it from lark-voucher and basic information。
		AppType:     protocol.InternalApp, //AppType only has two types: Independent Software Vendor App（ISVApp） or Internal App.
		AppSecret:   appSecret,            //get it from lark-voucher and basic information.
		VerifyToken: verifyToken,          //get it from lark-event subscriptions.
		EncryptKey:  encryptKey,           //get it from lark-event subscriptions.
	}

	return conf, nil
}
