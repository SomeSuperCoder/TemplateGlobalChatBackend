package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SomeSuperCoder/global-chat/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: routes.LoadRoutes(),
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	// ========== MongoDB ==========
	connectionString := "mongodb://localhost:27017"
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err = client.Ping(timeoutCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB nor connect: %w", err)
	}

	// database := client.Database("chat")

	// ========== HTTP server ==========
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
