package authentication

import (
	"context"
	"net/http"
	"time"

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
// client.InitDB(map[string]string{"addr": "127.0.0.1:6379"})
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

func (a *AuthMiniProgram) Login(ctx context.Context, code string, appInfo *AuthAppInfo, requestHost string) (map[string]*http.Cookie, error) {
	// check params
	if code == "" || appInfo == nil || appInfo.AppID == "" || appInfo.AppSecret == "" {
		return nil, common.ErrValidateParams.ErrorWithExtStr("miniProgram login input param is empty")
	}

	rsp, err := MiniProgramLoginValidate(code, appInfo.AppID, appInfo.AppSecret)
	if err != nil {
		return nil, err
	}

	sessionName := a.Manager.GenerateSessionKeyName(appInfo.AppID)

	authUser := GetAuthUserInfo(rsp)
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

func GetAuthUserInfo(rsp *protocol.MiniProgramLoginResponse) *AuthUserInfo {
	authUser := &AuthUserInfo{}

	authUser.Token.AccessToken = rsp.AccessToken
	authUser.Token.TokenType = rsp.TokenType
	authUser.Token.ExpiresIn = rsp.ExpiresIn
	authUser.Token.RefreshToken = rsp.RefreshToken

	authUser.User.TenantKey = rsp.TenantKey
	authUser.User.OpenID = rsp.OpenID
	authUser.User.EmployeeID = rsp.EmployeeID

	authUser.Extra = map[string]string{
		"uid":         rsp.UID,
		"union_id":    rsp.UnionID,
		"session_key": rsp.SessionKey,
	}

	return authUser
}
