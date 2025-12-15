package models

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Brand       string    `json:"brand"`
	Sku         string    `json:"sku"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	AddedBy     int64     `json:"-"`
}
