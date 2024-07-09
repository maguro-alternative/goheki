package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
)

func checkAuth(r *http.Request) bool {
	user, pass, ok := r.BasicAuth()
	return ok && subtle.ConstantTimeCompare([]byte(user), []byte(os.Getenv("BASIC_AUTH_USER_ID"))) == 1 &&
		subtle.ConstantTimeCompare([]byte(pass), []byte(os.Getenv("BASIC_AUTH_PASSWORD"))) == 1
}

func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !checkAuth(r) {
			w.Header().Add("WWW-Authenticate", `Basic realm="my private area"`)
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
