package middleware

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/utils"
)

func AuthMiddleware(next http.Handler, users map[string]utils.Login) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := utils.Authorize(r, users); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
