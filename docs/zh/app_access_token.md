# app_access_token 的获取  
获取 app_access_token 主要分为两部分。
1. 首先初始化 SDK 相关配置信息。  
2. 然后调用 SDK/auth/authorization.go 中的 GetAppAccessToken 函数获取 app_access_token 。  
  
## 初始化 SDK 相关配置信息  
请参考   
[SDK Init 文档](../../README.zh-cn.md)     
[SDK Init 示例代码](../../demo/sdk_init/sdk_init.go)  

## 获取 app_access_token  
通过 SDK/auth/authorization.go 中的 GetAppAccessToken 函数可以直接获取，示例代码如下：     
(在使用以下函数前需要先初始化 SDK 相关配置信息)  
```go
func GetAppAccessToken(ctx context.Context, appID string)(string,error){
   accessToken, err := auth.GetAppAccessToken(ctx, appID)
   if err != nil {
      return "", fmt.Errorf("get AppAccessToken failed[%v]", err)
   }
   return accessToken, nil
}
```
