package models

import "time"

type Item struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	ItemCode  int       `json:"item_code"`
	InStock   int       `json:"in_stock"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}
