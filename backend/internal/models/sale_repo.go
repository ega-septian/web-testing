package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SaleRepo records sale events. There's no real checkout flow yet, so these
// are seeded manually (or by whatever process eventually replaces that) —
// see CatalogHandler.RecordSale.
type SaleRepo struct {
	pool *pgxpool.Pool
}

func NewSaleRepo(pool *pgxpool.Pool) *SaleRepo {
	return &SaleRepo{pool: pool}
}

func (r *SaleRepo) Record(ctx context.Context, productID string, quantity int) (*Sale, error) {
	var s Sale
	err := r.pool.QueryRow(ctx,
		`INSERT INTO sales (product_id, quantity) VALUES ($1, $2)
		 RETURNING id, product_id, quantity, sold_at`,
		productID, quantity,
	).Scan(&s.ID, &s.ProductID, &s.Quantity, &s.SoldAt)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
