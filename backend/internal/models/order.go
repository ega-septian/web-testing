package models

// OrderItem is one line of an order — a snapshot of the product's name,
// brand, and price at the time of purchase, not a live reference (see
// order_items' migration comment: this is what keeps order history intact
// even if the product is later changed or deleted).
type OrderItem struct {
	ID          string   `json:"id"`
	ProductID   *string  `json:"product_id"`
	ProductName string   `json:"product_name"`
	Brand       string   `json:"brand"`
	Size        string   `json:"size"`
	UnitPrice   int      `json:"unit_price"`
	Quantity    int      `json:"quantity"`
	CreatedAt   JSONTime `json:"created_at"`
}

type Order struct {
	ID            string      `json:"id"`
	UserID        string      `json:"user_id"`
	RecipientName string      `json:"recipient_name"`
	Phone         string      `json:"phone"`
	Address       string      `json:"address"`
	TotalAmount   int         `json:"total_amount"`
	Status        string      `json:"status"`
	CreatedAt     JSONTime    `json:"created_at"`
	Items         []OrderItem `json:"items"`
}

// CheckoutItem is one requested line in a checkout request — just what the
// client can legitimately decide (which product+size, how many). Everything
// else (name, brand, current price, stock) is resolved and validated
// server-side in OrderRepo.Place, never trusted from the client.
type CheckoutItem struct {
	ProductID string
	Size      string
	Quantity  int
}
