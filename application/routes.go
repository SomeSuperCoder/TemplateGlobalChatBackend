package application

import (
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/handers"
	"github.com/SomeSuperCoder/global-chat/middleware"
)

func loadRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	mux.Handle("/auth/", loadAuthRoutes())

	return middleware.LoggerMiddleware(mux)
}

func loadAuthRoutes() http.Handler {
	authMux := http.NewServeMux()

	authHandler := &handers.AuthHandler{}

	authMux.HandleFunc("POST /register", authHandler.Register)
	authMux.HandleFunc("POST /login", authHandler.Login)
	authMux.HandleFunc("POST /logout", authHandler.Logout)

	return http.StripPrefix("/auth", authMux)
}
