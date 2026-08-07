package models

import "time"

type Asset struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Filename    string    `json:"filename"`
	URL         string    `json:"url"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
