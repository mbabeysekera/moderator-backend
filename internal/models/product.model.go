package models

import (
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
)

type Product struct {
	ID          int64                 `json:"id"`
	Title       string                `json:"title"`
	Brand       string                `json:"brand"`
	Category    enums.ProductCategory `json:"category"`
	Sku         string                `json:"sku"`
	InStock     int                   `json:"in_stock"`
	Description string                `json:"description"`
	Price       float64               `json:"price"`
	CreatedAt   time.Time             `json:"created_at"`
	AddedBy     int64                 `json:"-"`
	TenantID    int64                 `json:"tenant_id"`
}
