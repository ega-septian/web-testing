package models

import "time"

type DressStyle struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
