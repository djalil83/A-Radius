package auth

import "time"

type User struct {
	ID       string
	Username string
	Status   string
}

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
}
