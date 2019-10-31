# get app_access_token  
You can get it in two steps.    
1. Initialize SDK configuration.    
2. Use GetAppAccessToken function in SDK/auth/authorization.go to get app_access_token.    
  
## Initialize SDK configuration  
Please refer to   
[SDK Init doc](../../README.md)    
[SDK Init code](../../demo/sdk_init/sdk_init.go)    
  
## get app_access_token  
You can use GetAppAccessToken function in SDK/auth/authorization.go to get app_access_token easily.    
Of course, you should initialize SDK configuration first.    
```go
func GetAppAccessToken(ctx context.Context, appID string)(string,error){
   accessToken, err := auth.GetAppAccessToken(ctx, appID)
   if err != nil {
      return "", fmt.Errorf("get AppAccessToken failed[%v]", err)
   }
   return accessToken, nil
}
```
