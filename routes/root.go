package routes

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/middleware"
	"github.com/SomeSuperCoder/global-chat/utils"
)

var users = map[string]utils.Login{}

func LoadRoutes() http.Handler {
	rootMux := http.NewServeMux()

	rootMux.Handle("/", loadAuthRoutes())
	rootMux.Handle("/api/", middleware.AuthMiddleware(http.StripPrefix("/api", loadExampleMux()), users))

	return middleware.LoggerMiddleware(rootMux)
}
