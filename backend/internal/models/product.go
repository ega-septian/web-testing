package models

import "time"

type Product struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Section      string    `json:"section"`
	Icon         string    `json:"icon"`
	Rating       float64   `json:"rating"`
	Price        int       `json:"price"`
	OldPrice     *int      `json:"old_price"`
	DisplayOrder int       `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
