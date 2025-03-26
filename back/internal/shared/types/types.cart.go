package types

import "time"

type CartStore interface {
	CreateCart(userID int) error
	GetMyCartItems() (*[]Cartitem, error)
	AddItemToCart(productID int, userID int) (*Cartitem, error)
	RemoveItemFromCart(productID int, userID int) error
	GetTotal(userID int) (float64, error)
}

type CartService interface {
	CreateCart(userID int) error
	GetMyCartItems() *[]Cartitem
	AddItemToCart(productID int, userID int) *Cartitem
	RemoveItemFromCart(productID int, userID int)
	GetTotal(userID int) float64
}

type Cart struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Cartitem struct {
	CartID        int       `json:"cartId"`
	ProductID     int       `json:"productId"`
	Quantity      int       `json:"quantity"`
	PriceAtAdding float64   `json:"priceAtAdding"`
	AddedAt       time.Time `json:"addedAt"`
}
