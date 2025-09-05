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

	mux.HandleFunc("POST /health", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}, db))
	mux.Handle("/auth/", loadAuthRoutes(db))
	mux.Handle("/messages/", loadMessageRoutes(db))

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
	authMux.HandleFunc("POST /logout", middleware.AuthMiddleware(authHandler.Logout, db))

	return http.StripPrefix("/auth", authMux)
}

func loadMessageRoutes(db *mongo.Database) http.Handler {
	messageMux := http.NewServeMux()
	messageHandler := &handers.MessageHandler{
		Repo: repository.MessageRepo{
			Database: db,
		},
	}

	messageMux.HandleFunc("GET /", middleware.AuthMiddleware(messageHandler.CreateMessage, db))
	messageMux.HandleFunc("CREATE /", middleware.AuthMiddleware(messageHandler.CreateMessage, db))
	messageMux.HandleFunc("UPDATE /{id}", middleware.AuthMiddleware(messageHandler.CreateMessage, db))
	messageMux.HandleFunc("DELETE /{id}", middleware.AuthMiddleware(messageHandler.CreateMessage, db))

	return http.StripPrefix("/messages", messageMux)
}
