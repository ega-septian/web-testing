package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DressStyleRepo struct {
	pool *pgxpool.Pool
}

func NewDressStyleRepo(pool *pgxpool.Pool) *DressStyleRepo {
	return &DressStyleRepo{pool: pool}
}

func (r *DressStyleRepo) List(ctx context.Context) ([]DressStyle, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, display_order, created_at FROM dress_styles ORDER BY display_order ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	styles := []DressStyle{}
	for rows.Next() {
		var s DressStyle
		if err := rows.Scan(&s.ID, &s.Name, &s.DisplayOrder, &s.CreatedAt); err != nil {
			return nil, err
		}
		styles = append(styles, s)
	}

	return styles, rows.Err()
}
