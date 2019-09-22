package message_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/message"
)

var (
	once      sync.Once
	appConf   *appconfig.AppConfig
	tenantKey string
	chatID    string
	openID    string
	userID    string
	imageKey  string
)

func InitTestParams() {
	once.Do(func() {
		appConf = &appconfig.AppConfig{
			AppID:       os.Getenv("appid"),
			AppSecret:   os.Getenv("appsecret"),
			VerifyToken: os.Getenv("verifytoken"),
			EncryptKey:  os.Getenv("encryptkey"),
			AppType:     os.Getenv("apptype"),
		}

		tenantKey = os.Getenv("tenantkey")
		chatID = os.Getenv("chatid")
		openID = os.Getenv("openid")
		userID = os.Getenv("userid")

		imageKey = os.Getenv("imagekey")

		appconfig.Init(*appConf)
	})
}

func TestGetImageKey(t *testing.T) {
	c := context.Background()
	InitTestParams()

	// by path
	path := "/tmp/test.png"
	key, err := message.GetImageKey(c, tenantKey, appConf.AppID, "", path)
	if err != nil {
		t.Errorf("GetImageKeyByPath failed: %v", err)
	} else {
		t.Logf("GetImageKeyByPath: %+v", key)
	}

	// by url
	url := "http://a.hiphotos.baidu.com/image/pic/item/838ba61ea8d3fd1fc9c7b6853a4e251f94ca5f46.jpg"
	key, err = message.GetImageKey(c, tenantKey, appConf.AppID, url, "")
	if err != nil {
		t.Errorf("GetImageKeyByURL failed: %v", err)
	} else {
		t.Logf("GetImageKeyByURL: %+v", key)
	}
}
