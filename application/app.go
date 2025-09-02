package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type App struct {
	router http.Handler
	client *mongo.Client
	db     *mongo.Database
}

func New() *App {
	app := &App{}

	return app
}

func (a *App) Start(ctx context.Context) error {
	var err error
	// ========== MongoDB ==========
	connectionString := "mongodb://localhost:27017"
	a.client, err = mongo.Connect(options.Client().ApplyURI(connectionString))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err = a.client.Ping(timeoutCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB nor connect: %w", err)
	}

	// Get the project database
	a.db = a.client.Database("chat")

	// ========== Load Routes ==========
	a.router = loadRoutes(a.db)

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
