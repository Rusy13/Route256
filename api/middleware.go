package postgresql

import (
	"context"
	"fmt"
	"net/http"
)

const UserContextKey = "user"

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		if ok && VerifyUserPass(user, pass) {
			newctx := context.WithValue(req.Context(), UserContextKey, user)
			next.ServeHTTP(w, req.WithContext(newctx))
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		fmt.Println("SUCCESS")
	})
}
