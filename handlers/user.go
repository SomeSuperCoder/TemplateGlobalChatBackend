package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SomeSuperCoder/global-chat/middleware"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserHandler struct {
	Repo repository.UserRepo
}

// ==============================================================
// ================ Auth-related handlers =======================
// ==============================================================

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username") // Allowed
	password := r.FormValue("password") // Allowed

	// Check username and password length
	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "Invalid username/password", http.StatusNotAcceptable)
		return
	}

	// Make sure such user does not already exist
	doesExist := h.Repo.DoesExist(r.Context(), username)
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

	err := h.Repo.CreateUser(r.Context(), newUser)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User registered successfully!")

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username") // Allowed
	password := r.FormValue("password") // Allowed

	// Get the user
	user, err := h.Repo.GetUserByUsername(r.Context(), username)

	// Check if user exists
	if errors.Is(err, repository.ErrEntryNotFound) {
		http.Error(w, "Wrong username or password", http.StatusUnauthorized)
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
		Path:     "/",
	})

	// Set CSRF token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  expires,
		HttpOnly: false,
		Path:     "/",
	})

	// Store token in DB
	newSession := models.UserSession{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		CratedAt:     time.Now(),
	}
	err = h.Repo.AddLoginSession(r.Context(), username, newSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Login successful!")
}

// This functions needs to be wrapped with an auth middleware
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userAuth := middleware.ExtractUserAuth(r)

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
	err := h.Repo.FinalizeSession(r.Context(), userAuth.Username, sessionToken.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintln(w, "Logged out successfully!")
}

// ==============================================================
// ================ Non-auth-related handlers ===================
// ==============================================================
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	parsedUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID provided", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetUserByID(r.Context(), parsedUserID)
	if err != nil {
		if err == repository.ErrEntryNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serializedUser, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(serializedUser))
}
