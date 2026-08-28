package models

type User struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"-"`
	CreatedAt    JSONTime `json:"created_at"`
}
