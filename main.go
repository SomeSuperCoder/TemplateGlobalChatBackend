package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SomeSuperCoder/go-auth/internal/middleware"
	"github.com/SomeSuperCoder/go-auth/internal/utils"
)

var users = map[string]utils.Login{}

func main() {
	rootMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	// Auth routes
	rootMux.HandleFunc("POST /register", register)
	rootMux.HandleFunc("POST /login", login)
	rootMux.Handle("POST /logout", middleware.AuthMiddleware(http.HandlerFunc(logout), users))

	// API routes
	apiMux.HandleFunc("POST /protected", protected)
	rootMux.Handle("/api/", middleware.AuthMiddleware(http.StripPrefix("/api", apiMux), users))

	http.ListenAndServe(":8080", rootMux)
}

func register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "Invalid username/password", http.StatusNotAcceptable)
		return
	}

	if _, ok := users[username]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	users[username] = utils.Login{
		HashedPassword: hashedPassword,
	}

	fmt.Fprintln(w, "User registered successfully!")

}

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, ok := users[username]
	if !ok || !utils.CheckPasswordhash(password, user.HashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)
	expires := time.Now().Add(24 * time.Hour)

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
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user

	fmt.Fprintln(w, "Login successful!")
}

func logout(w http.ResponseWriter, r *http.Request) {
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

	username := r.FormValue("username")
	user, _ := users[username]
	user.SessionToken = ""
	user.CSRFToken = ""
	users[username] = user

	fmt.Fprintln(w, "Logged out successfully!")
}

func protected(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	fmt.Fprintf(w, "Welcome, %s!\n", username)
}
