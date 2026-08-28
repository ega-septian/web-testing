package models

type DressStyle struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	DisplayOrder int      `json:"-"`
	CreatedAt    JSONTime `json:"created_at"`
}
