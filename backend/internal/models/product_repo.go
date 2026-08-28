package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidSection = errors.New("invalid product section")

var validProductSections = map[string]bool{
	"new_arrivals": true,
	"top_selling":  true,
}

type ProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{pool: pool}
}

// ListBySection returns products for one homepage section, ordered for display.
// Returns ErrInvalidSection if section isn't one of the known values, rather than
// silently returning an empty list for a typo'd query param.
func (r *ProductRepo) ListBySection(ctx context.Context, section string) ([]Product, error) {
	if !validProductSections[section] {
		return nil, ErrInvalidSection
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, section, icon, rating, price, old_price, display_order, created_at
		 FROM products WHERE section = $1 ORDER BY display_order ASC`,
		section,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Section, &p.Icon, &p.Rating, &p.Price, &p.OldPrice, &p.DisplayOrder, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, rows.Err()
}
