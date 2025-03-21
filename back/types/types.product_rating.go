package types

import "time"

type ProductRatingStore interface {
	CreateRating(*CreateProductRatingPayload, int) (*ProductRating, error)
	GetRatingsByProduct(int) ([]*ProductRating, error)
	GetRatingsByUser(int) ([]*ProductRating, error)
	GetRating(int) (*ProductRating, error)
	UpdateRating(int, *UpdateProductRatingPayload) (*ProductRating, error)
	DeleteRating(int) error
}

type ProductRating struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	ProductID int       `json:"productId"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateProductRatingPayload struct {
	ProductID int    `json:"productId" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment" validate:"max=1000"`
}

type UpdateProductRatingPayload struct {
	Rating  int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment string `json:"comment,omitempty" validate:"omitempty,max=1000"`
}
