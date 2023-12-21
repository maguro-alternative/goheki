package cookie

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/maguro-alternative/goheki/configs/envconfig"
)

// セッションの定義
var (
	Store *sessions.CookieStore
)

func init() {
	env, err := envconfig.NewEnv()
	if err != nil {
		panic(err)
	}
	Store = sessions.NewCookieStore([]byte(env.SessionsSecret))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	// ドメインが設定されている場合はセット
	if env.CookieDomain != "" {
		Store.Options.Domain = env.CookieDomain
	}
}