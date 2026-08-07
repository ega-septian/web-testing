package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAssetNotFound = errors.New("asset not found")

type AssetRepo struct {
	pool *pgxpool.Pool
}

func NewAssetRepo(pool *pgxpool.Pool) *AssetRepo {
	return &AssetRepo{pool: pool}
}

// Upsert creates or replaces the asset stored under the given key.
func (r *AssetRepo) Upsert(ctx context.Context, key, filename, url, contentType string, sizeBytes int64) (*Asset, error) {
	var a Asset
	err := r.pool.QueryRow(ctx,
		`INSERT INTO assets (key, filename, url, content_type, size_bytes)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (key) DO UPDATE
		   SET filename = EXCLUDED.filename,
		       url = EXCLUDED.url,
		       content_type = EXCLUDED.content_type,
		       size_bytes = EXCLUDED.size_bytes,
		       updated_at = now()
		 RETURNING id, key, filename, url, content_type, size_bytes, created_at, updated_at`,
		key, filename, url, contentType, sizeBytes,
	).Scan(&a.ID, &a.Key, &a.Filename, &a.URL, &a.ContentType, &a.SizeBytes, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *AssetRepo) FindByKey(ctx context.Context, key string) (*Asset, error) {
	var a Asset
	err := r.pool.QueryRow(ctx,
		`SELECT id, key, filename, url, content_type, size_bytes, created_at, updated_at
		 FROM assets WHERE key = $1`,
		key,
	).Scan(&a.ID, &a.Key, &a.Filename, &a.URL, &a.ContentType, &a.SizeBytes, &a.CreatedAt, &a.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAssetNotFound
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *AssetRepo) List(ctx context.Context) ([]Asset, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, key, filename, url, content_type, size_bytes, created_at, updated_at
		 FROM assets ORDER BY key ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assets := []Asset{}
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.Key, &a.Filename, &a.URL, &a.ContentType, &a.SizeBytes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	return assets, rows.Err()
}

// Delete removes the asset row and returns its filename so the caller can remove the file on disk.
func (r *AssetRepo) Delete(ctx context.Context, key string) (string, error) {
	var filename string
	err := r.pool.QueryRow(ctx,
		`DELETE FROM assets WHERE key = $1 RETURNING filename`,
		key,
	).Scan(&filename)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrAssetNotFound
	}
	if err != nil {
		return "", err
	}

	return filename, nil
}
