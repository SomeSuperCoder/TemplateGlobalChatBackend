package handers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
)

var users = map[string]utils.Login{}

type AuthHandler struct {
	Repo repository.UserRepo
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check username and password length
	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "Invalid username/password", http.StatusNotAcceptable)
		return
	}

	// Make sure such user does not already exist
	doesExist := h.Repo.DoesExist(context.TODO(), username)
	if doesExist {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Create new user
	hashedPassword, _ := utils.HashPassword(password)
	newUser := &models.User{
		Username:       username,
		HashedPassword: hashedPassword,
		Sessions:       []models.UserSession{},
		CratedAt:       time.Now(),
	}

	h.Repo.CreateUser(context.TODO(), newUser)

	fmt.Fprintln(w, "User registered successfully!")

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Get the user
	user, err := h.Repo.GetUser(context.TODO(), username)

	// Check if user exists
	if errors.Is(err, repository.ErrUserNotFound) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check for any other errors
	if err != nil {
		fmt.Println("Any other error!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify password
	if !utils.CheckPasswordhash(password, user.HashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate tokens and expire date
	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)
	expires := time.Now().Add(7 * 24 * time.Hour) // 1 week

	// Set token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expires,
		HttpOnly: true,
	})

	// Set CSRF token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  expires,
		HttpOnly: false,
	})

	// Store token in DB
	newSession := models.UserSession{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		CratedAt:     time.Now(),
	}
	err = h.Repo.AddLoginSession(context.TODO(), username, newSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Login successful!")
}

// This functions needs to be wrapped with an auth middleware
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	// Reset cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	// Get session token
	sessionToken, _ := r.Cookie("session_token")
	// Remove session from database
	err := h.Repo.FinalizeSession(context.TODO(), username, sessionToken.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintln(w, "Logged out successfully!")
}
