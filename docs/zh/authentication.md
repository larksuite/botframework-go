# 身份认证  
`SDK/authentication/authentication.go`: `Authentication` 定义身份认证接口，`SessionManager` 定义用户 session 存取接口。  
  
`SDK/authentication/default_session_manager.go`: SDK 提供默认的 seesion 存取实现，通过`NewDefaultSessionManager`函数获取，代码示例为  
```golang
client := &common.DefaultRedisClient{}
client.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
manager := authentication.NewDefaultSessionManager("DojK2hs*790(", client)
```
  
`SDK/authentication/mini_program.go`: SDK 提供了小程序身份认证相关接口  
`SDK/authentication/open_sso.go`: SDK 提供了 open sso 身份认证相关接口  
`SDK/authentication/auth_mini_program.go`: SDK 已经实现小程序的登录、登录态校验、登出操作  
  
## 小程序身份认证示例代码  
初始化  
```golang
var minaAuth *authentication.AuthMiniProgram

func InitMinaAuth() error {
	client := &common.DefaultRedisClient{}
	err := client.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
	if err != nil {
		common.Logger(context.TODO()).Errorf("init redis error[%v]", err)
		return fmt.Errorf("init redis error[%v]", err)
	}
	common.Logger(context.TODO()).Infof("init redis succss")

	minaAuth = authentication.NewAuthMiniProgram(authentication.NewDefaultSessionManager("DojK2hs*790k", client), time.Hour*24*7)

	return nil
}
```
  
登录示例代码  
```golang
func MinaLogin(c *gin.Context) {
	appID := "cli_9d1ad8ed77f69108"

	params, err := url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		common.Logger(c).Errorf("minaLogin: getParams[%v]", err)
		c.JSON(500, gin.H{"codemsg": common.ErrMinaCodeGetParams.String()})
		return
	}

	code := params.Get("code")
	if code == "" {
		common.Logger(c).Errorf("minaLogin: code is empty")
		c.JSON(500, gin.H{"codemsg": common.ErrMinaCodeGetParams.String()})
		return
	}

	host := c.Request.Host

	mapCookie, err := minaAuth.Login(c, code, appID, host)
	if err != nil {
		common.Logger(c).Errorf("minaLogin: codeToSession[%v]", err)
		c.JSON(500, gin.H{"codemsg": fmt.Sprintf("%v", err)})
		return
	}

	for _, v := range mapCookie {
		http.SetCookie(c.Writer, v)
		common.Logger(c).Infof("minaLogin: code[%s]host[%s]cookie[%+v]", code, host, v)
	}

	c.JSON(200, gin.H{"code": 0, "msg": ""})
	return
}
```
  
校验登录状态示例代码  
```golang
func MinaCheckAuth(c *gin.Context) {
	appID := "cli_9d1ad8ed77f69108"

	sessionKeyName := minaAuth.GetSessionManager().GenerateSessionKeyName(appID)
	sessionKey, err := c.Cookie(sessionKeyName)
	if err != nil {
		common.Logger(c).Errorf("MinaCheckAuth: getSessionKey error[%v]", err)
		c.JSON(500, gin.H{"codemsg": common.ErrMinaSessionInvalid.String()})
		return
	}
	common.Logger(c).Infof("MinaCheckAuth: sessionKeyName[%s]sessionKey[%s]", sessionKeyName, sessionKey)

	err = minaAuth.Auth(c, sessionKey)
	if err != nil {
		common.Logger(c).Errorf("MinaCheckAuth: Auth error[%v]", err)
		c.JSON(500, gin.H{"codemsg": common.ErrMinaSessionInvalid.String()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": ""})
	return
}
```
  
登出示例代码  
```golang
func MinaLogout(c *gin.Context) {
	appID := "cli_9d1ad8ed77f69108"

	mapCookie, err := minaAuth.Logout(c, appID, c.Request.Host)
	if err != nil {
		common.Logger(c).Errorf("minaLogout: logout error[%v]", err)
		c.JSON(500, gin.H{"codemsg": common.ErrMinaSessionInvalid.String()})
		return
	}

	for _, v := range mapCookie {
		http.SetCookie(c.Writer, v)
		common.Logger(c).Infof("minaLogout: cookie[%+v]", v)
	}

	c.JSON(200, gin.H{"code": 0, "msg": ""})
	return
}
```
