package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductSizeRepo struct {
	pool *pgxpool.Pool
}

func NewProductSizeRepo(pool *pgxpool.Pool) *ProductSizeRepo {
	return &ProductSizeRepo{pool: pool}
}

func (r *ProductSizeRepo) ListByProduct(ctx context.Context, productID string) ([]ProductSize, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT size, stock FROM product_sizes WHERE product_id = $1 ORDER BY size ASC`,
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sizes := []ProductSize{}
	for rows.Next() {
		var s ProductSize
		if err := rows.Scan(&s.Size, &s.Stock); err != nil {
			return nil, err
		}
		sizes = append(sizes, s)
	}

	return sizes, rows.Err()
}
