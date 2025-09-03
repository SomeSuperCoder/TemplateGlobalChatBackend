package middleware

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func AuthMiddleware(next http.HandlerFunc, db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := utils.Authorize(r, db); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
