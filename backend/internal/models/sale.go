package models

type Sale struct {
	ID        string   `json:"id"`
	ProductID string   `json:"product_id"`
	Quantity  int      `json:"quantity"`
	SoldAt    JSONTime `json:"sold_at"`
}
