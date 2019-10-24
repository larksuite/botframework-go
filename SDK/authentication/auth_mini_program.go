package authentication

import (
	"context"
	"net/http"
	"time"

	"github.com/larksuite/botframework-go/SDK/auth"
	"github.com/larksuite/botframework-go/SDK/common"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

type AuthMiniProgram struct {
	Manager         SessionManager
	ValidPeriod     time.Duration
	AuthCookieLevel CookieDomainLevel
}

// NewAuthMiniProgram demo:
// client := &common.DefaultRedisClient{}
// err := client.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
// if err != nil {
// 	return fmt.Errorf("init redis error[%v]", err)
// }
// manager := authentication.NewDefaultSessionManager("DojK2hs*790(", client)
// minaAuth := authentication.NewAuthMiniProgram(manager, time.Hour*24*7)
func NewAuthMiniProgram(manager SessionManager, validPeriod time.Duration) *AuthMiniProgram {
	mina := &AuthMiniProgram{}
	mina.Init(manager, validPeriod)

	return mina
}

func (a *AuthMiniProgram) Init(manager SessionManager, validPeriod time.Duration) {
	a.Manager = manager
	a.ValidPeriod = validPeriod

	a.SetCookieDomainLevel(DomainLevelZero) // DomainLevelZero is default auth cookie level
}

func (a *AuthMiniProgram) Login(ctx context.Context, code string, appID string, requestHost string) (map[string]*http.Cookie, error) {
	// check params
	if code == "" || appID == "" {
		return nil, common.ErrValidateParams.ErrorWithExtStr("miniProgram login input param is empty")
	}

	// get app_access_token
	appAccessToken, err := auth.GetAppAccessToken(ctx, appID)
	if err != nil {
		return nil, err
	}

	rsp, err := MiniProgramValidateByAppToken(code, appAccessToken)
	if err != nil {
		return nil, err
	}

	sessionName := a.Manager.GenerateSessionKeyName(appID)

	authUser := TransMPAuthUser(rsp)

	sessionKey, err := a.Manager.SetAuthUserInfo(authUser, a.GetValidPeriod())
	if err != nil {
		return nil, common.ErrMinaSetAuth.ErrorWithExtErr(err)
	}

	// cookie
	mapCookie := make(map[string]*http.Cookie)
	mapCookie[sessionName] = &http.Cookie{
		Name:    sessionName,
		Value:   sessionKey,
		Expires: time.Now().Add(a.GetValidPeriod()),
		MaxAge:  int(a.GetValidPeriod().Seconds()),
		Path:    "/",
		Domain:  GetAuthCookieDomain(requestHost, a.GetCookieDomainLevel()),
	}

	return mapCookie, nil
}

func (a *AuthMiniProgram) Auth(ctx context.Context, sessionKey string) error {
	// check params
	if sessionKey == "" {
		return common.ErrAuthParams.ErrorWithExtStr("miniProgram Auth input param is empty")
	}

	authUser, err := a.Manager.GetAuthUserInfo(sessionKey)
	if err != nil {
		return common.ErrMinaGetAuth.ErrorWithExtErr(err)
	}

	if len(authUser.User.OpenID) == 0 {
		return common.ErrMinaSessionInvalid.Error()
	}

	return nil
}

func (a *AuthMiniProgram) Logout(ctx context.Context, appID string, requestHost string) (map[string]*http.Cookie, error) {
	// check params
	if appID == "" {
		return nil, common.ErrLogoutParams.ErrorWithExtStr("miniProgram Auth input param is empty")
	}

	sessionName := a.Manager.GenerateSessionKeyName(appID)

	mapCookie := make(map[string]*http.Cookie)
	mapCookie[sessionName] = &http.Cookie{
		Name:    sessionName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
		MaxAge:  -1,
		Path:    "/",
		Domain:  GetAuthCookieDomain(requestHost, a.GetCookieDomainLevel()),
	}

	return mapCookie, nil

}

func (a *AuthMiniProgram) SetCookieDomainLevel(cookieLevel CookieDomainLevel) {
	a.AuthCookieLevel = cookieLevel
}

func (a *AuthMiniProgram) GetCookieDomainLevel() CookieDomainLevel {
	return a.AuthCookieLevel
}

func (a *AuthMiniProgram) GetValidPeriod() time.Duration {
	return a.ValidPeriod
}

func (a *AuthMiniProgram) GetSessionManager() SessionManager {
	return a.Manager
}

func TransMPAuthUser(rsp *protocol.MiniProgramLoginByAppTokenResponse) *AuthUserInfo {
	authUser := &AuthUserInfo{}

	authUser.Token.AccessToken = rsp.Data.AccessToken
	authUser.Token.TokenType = rsp.Data.TokenType
	authUser.Token.ExpiresIn = rsp.Data.ExpiresIn
	authUser.Token.RefreshToken = rsp.Data.RefreshToken

	authUser.User.TenantKey = rsp.Data.TenantKey
	authUser.User.OpenID = rsp.Data.OpenID

	authUser.Extra = map[string]string{
		"union_id":    rsp.Data.UnionID,
		"session_key": rsp.Data.SessionKey,
	}

	return authUser
}
