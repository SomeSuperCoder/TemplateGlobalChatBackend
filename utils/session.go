package utils

import (
	"errors"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request, db *mongo.Database) (*repository.UserAuth, error) {
	// Init repo
	repo := repository.UserRepo{
		Database: db,
	}

	// Get the session token from the cookie
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" {
		return nil, err
	}

	// Get the CSRF token from the headers
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf == "" {
		return nil, AuthError
	}

	// Verify
	userAuth, err := repo.AuthCheck(r.Context(), st.Value, csrf)
	if err != nil {
		return nil, err
	}

	return userAuth, nil
}
