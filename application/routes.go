package application

import (
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/handers"
	"github.com/SomeSuperCoder/global-chat/middleware"
	"github.com/SomeSuperCoder/global-chat/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func loadRoutes(db *mongo.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	mux.Handle("/auth/", loadAuthRoutes(db))

	return middleware.LoggerMiddleware(mux)
}

func loadAuthRoutes(db *mongo.Database) http.Handler {
	authMux := http.NewServeMux()
	authHandler := &handers.AuthHandler{
		Repo: repository.UserRepo{
			Database: db,
		},
	}

	authMux.HandleFunc("POST /register", authHandler.Register)
	authMux.HandleFunc("POST /login", authHandler.Login)
	authMux.HandleFunc("POST /logout", authHandler.Logout)

	return http.StripPrefix("/auth", authMux)
}
