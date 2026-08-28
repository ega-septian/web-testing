package models

type ProductImage struct {
	ID           string   `json:"id"`
	ProductID    string   `json:"-"`
	URL          string   `json:"url"`
	DisplayOrder int      `json:"display_order"`
	CreatedAt    JSONTime `json:"created_at"`
}
