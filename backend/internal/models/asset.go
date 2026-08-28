package models

type Asset struct {
	ID          string   `json:"id"`
	Key         string   `json:"key"`
	Filename    string   `json:"filename"`
	URL         string   `json:"url"`
	ContentType string   `json:"content_type"`
	SizeBytes   int64    `json:"size_bytes"`
	CreatedAt   JSONTime `json:"created_at"`
	UpdatedAt   JSONTime `json:"updated_at"`
}
