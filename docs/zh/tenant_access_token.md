# tenant_access_token 的获取和管理  
获取 tenant_access_token 主要分为两部分。  
1. 首先初始化 SDK 相关配置信息。  
2. 然后调用 SDK/auth/authorization.go 中的 GetTenantAccessToken 函数获取 tenant_access_token 。  
  
## 初始化 SDK 相关配置信息  
请参考[SDK Init](../../README.zh-cn.md).  

## 获取 tenant_access_token  
gin 框架生成代码:    
使用自动生成的代码无需特地获取 tenant_access_token ，在需要使用 tenant_access_token 的地方，如发送文本消息时使用的 SendTextMessage 函数，参数中有 tenantKey 一项，输入正确的 tenantKey 即可。    
SDK 会自动使用 tenantKey 去获取 tenant_access_token ，而无需开发者手动获取。    
SendTextMessage 函数如下：    
```go
func SendTextMessage(ctx context.Context, tenantKey, appID string,
   user *protocol.UserInfo, rootID string,
   text string) (*protocol.SendMsgResponse, error) {
   return sendMsg(ctx, tenantKey, appID, protocol.NewTextMsgReq(user, rootID, text))
}
```
  
手动编写代码:    
通过 SDK/auth/authorization.go 中的 GetTenantAccessToken 函数可以直接获取，示例代码如下:    
(在使用以下函数前需要先初始化 SDK 相关配置信息)  
```go
func GetTenantAccessToken(ctx context.Context, tenantKey, appID string)(string,error){
   accessToken, err := auth.GetTenantAccessToken(ctx, tenantKey, appID)
   if err != nil {
      return "", fmt.Errorf("get TenantAccessToken failed[%v]", err)
   }
   return accessToken, nil
}
```
