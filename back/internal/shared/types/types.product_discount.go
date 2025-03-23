package types

import "time"

type ProductDiscountStore interface {
	CreateDiscount(*CreateProductDiscountPayload) (*ProductDiscount, error)
	UpdateDiscount(int, *UpdateProductDiscountPayload) (*ProductDiscount, error)
	DeleteDiscount(int) error
	GetDiscoutsByID(int) (*ProductDiscount, error)
	GetDiscountsByProduct(int) ([]*ProductDiscount, error)
	GetActiveDiscounts(int) ([]*ProductDiscount, error)
	GetDiscountsByDateRange(int, time.Time, time.Time) ([]*ProductDiscount, error)
}

type ProductDiscount struct {
	ID              int       `json:"id"`
	ProductID       int       `json:"productId"`
	DiscountPercent float64   `json:"discountPercent"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	CreatedAt       time.Time `json:"createdAt"`
}

type CreateProductDiscountPayload struct {
	ProductID       int       `json:"productId" validate:"required"`
	DiscountPercent float64   `json:"discountPercent" validate:"required"`
	StartDate       time.Time `json:"startDate" validate:"required"`
	EndDate         time.Time `json:"endDate" validate:"required"`
}

type UpdateProductDiscountPayload struct {
	DiscountPercent *float64   `json:"discountPercent,omitempty" validate:"omitempty"`
	StartDate       *time.Time `json:"startDate,omitempty" validate:"omitempty"`
	EndDate         *time.Time `json:"endDate,omitempty" validate:"omitempty"`
}
