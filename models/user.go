package models

type UserSession struct {
	sessionToken string
	csrfToken    string
}

type User struct {
	username string
	password string
	sessions []UserSession
}
