package models

import "time"

type Item struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	TenantID  int64     `json:"tenant_id"`
}
