package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MaxProductImages caps the gallery at 3 photos per product (main + 2), matching
// the approved product-detail design — not an arbitrary technical limit.
const MaxProductImages = 3

type ProductImageRepo struct {
	pool *pgxpool.Pool
}

func NewProductImageRepo(pool *pgxpool.Pool) *ProductImageRepo {
	return &ProductImageRepo{pool: pool}
}

func (r *ProductImageRepo) ListByProduct(ctx context.Context, productID string) ([]ProductImage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, product_id, url, display_order, created_at
		 FROM product_images WHERE product_id = $1 ORDER BY display_order ASC`,
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []ProductImage{}
	for rows.Next() {
		var img ProductImage
		if err := rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.DisplayOrder, &img.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, rows.Err()
}

// Add appends an image to the end of the product's gallery. Caller is
// responsible for enforcing MaxProductImages before calling this.
func (r *ProductImageRepo) Add(ctx context.Context, productID, url string) (*ProductImage, error) {
	var img ProductImage
	err := r.pool.QueryRow(ctx,
		`INSERT INTO product_images (product_id, url, display_order)
		 VALUES ($1, $2, (SELECT COALESCE(MAX(display_order), 0) + 1 FROM product_images WHERE product_id = $1))
		 RETURNING id, product_id, url, display_order, created_at`,
		productID, url,
	).Scan(&img.ID, &img.ProductID, &img.URL, &img.DisplayOrder, &img.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &img, nil
}
