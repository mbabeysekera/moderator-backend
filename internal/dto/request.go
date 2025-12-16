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
	ItemCode int    `json:"item_code"`
	InStock  int    `json:"in_stock"`
	ImageURL string `json:"image_url"`
}

type ProductsWithItemsRequest struct {
	Title       string              `json:"title"`
	Brand       string              `json:"brand"`
	Sku         string              `json:"sku"`
	Description string              `json:"description"`
	Items       []ItemCreateRequest `json:"items"`
}
