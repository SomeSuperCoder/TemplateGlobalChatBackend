package utils

import (
	"errors"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request, db *mongo.Database) error {
	// Init repo
	repo := repository.UserRepo{
		Database: db,
	}

	// Get uername
	username := r.FormValue("username")

	// Get the session token from the cookie
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" {
		return AuthError
	}

	// Get the CSRF token from the headers
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf == "" {
		return AuthError
	}

	// Verify
	if !repo.AuthCheck(r.Context(), username, st.Value, csrf) {
		return AuthError
	}

	return nil
}
