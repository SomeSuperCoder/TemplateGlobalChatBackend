package middleware

import (
	"context"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const UserAuthKey = "userAuth"

func ExtractUserAuth(r *http.Request) *repository.UserAuth {
	userAuth, ok := r.Context().Value(UserAuthKey).(*repository.UserAuth)
	if !ok {
		panic("Failed to extract user auth data")
	}

	return userAuth
}

func AuthMiddleware(next http.HandlerFunc, db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAuth, err := utils.Authorize(r, db)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserAuthKey, userAuth)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
