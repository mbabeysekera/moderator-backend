package dto

import (
	enums "coolbreez.lk/moderator/internal/constants"
)

type EventRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Target      string          `json:"target"`
	Command     string          `json:"command"`
	Status      enums.EventCode `json:"status"`
}

type UserCreateRequest struct {
	MobileNo string `json:"mobile_no" binding:"required,min=10,max=15,numeric"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"max=128"`
}

type UserLoginRequest struct {
	MobileNo string `json:"mobile_no" binding:"required,min=10,max=15,numeric"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserUpdateDetails struct {
	ID       int64  `json:"id" binding:"required"`
	MobileNo string `json:"mobile_no" binding:"required,min=10,max=15,numeric"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name" binding:"max=128"`
}

type ItemCreateRequest struct {
	ImageURL string `json:"image_url"`
}

type ProductsWithItemsRequest struct {
	AppID       int64                 `json:"app_id" binding:"required"`
	Title       string                `json:"title" binding:"required,max=24"`
	Brand       string                `json:"brand" binding:"required,max=24"`
	Category    enums.ProductCategory `json:"category" binding:"required"`
	Sku         string                `json:"sku"`
	Description string                `json:"description" binding:"required,max=64"`
	Price       float64               `json:"price" binding:"required"`
	InStock     int                   `json:"in_stock"`
	Items       []ItemCreateRequest   `json:"items" binding:"required"`
}

type ProductDetailsUpdateRequest struct {
	ID      int64    `json:"id" binding:"required"`
	AppID   int64    `json:"app_id" binding:"required"`
	Price   *float64 `json:"price,omitempty"`
	InStock *int     `json:"in_stock,omitempty"`
}
