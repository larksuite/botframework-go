# get and manage tenant_access_token  
You can get it in two steps.  
1. Initialize SDK configuration.  
2. Use GetTenantAccessToken function in SDK/auth/authorization.go to get tenant_access_token.  

## Initialize SDK configuration  
Please refer to   
[SDK Init doc](../../README.md).    
[SDK Init code](../../demo/sdk_init/sdk_init.go)    

## get tenant_access_token  
generated code:     
If you use generated code,  you just need to provide param tenantKey, and SDK will use tenantKey to get tenant_access_token automatically.    
For example, when you use SendTextMessage function in SDK/message/message.go.    
```go
func SendTextMessage(ctx context.Context, tenantKey, appID string,
   user *protocol.UserInfo, rootID string,
   text string) (*protocol.SendMsgResponse, error) {
   return sendMsg(ctx, tenantKey, appID, protocol.NewTextMsgReq(user, rootID, text))
}
```

your own code:    
You can use GetTenantAccessToken function in SDK/auth/authorization.go to get tenant_access_token easily.    
Of course, you should initialize SDK configuration first.    
```go
func GetTenantAccessToken(ctx context.Context, tenantKey, appID string)(string,error){
   accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
   if err != nil {
      return "", fmt.Errorf("get TenantAccessToken failed[%v]", err)
   }
   return accessToken, nil
}
```
