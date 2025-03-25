package types

import (
	"time"
)

type UserFavoriteStore interface {
	AddFavorite(userID int, productID int) (*UserFavorite, error)
	RemoveFavorite(userID int, productID int) error
	GetUserFavorite(userID int) (*[]*UserFavorite, error)
	GetFavorite(userID int, productID int) (*UserFavorite, error)
	IsFavorited(userID int, productID int) (bool, error)
}

type UserFavoriteService interface {
	AddFavorite(userID int, productID int) *UserFavorite
	RemoveFavorite(userID int, productID int)
	GetUserFavorite(userID int) *[]*UserFavorite
	IsFavorited(userID int, productID int) bool
}

type UserFavorite struct {
	UserID    int       `json:"id"`
	ProductId int       `json:"productId"`
	AddedAt   time.Time `json:"addedAt"`
}
