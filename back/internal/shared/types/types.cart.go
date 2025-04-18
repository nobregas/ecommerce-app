package types

import "time"

type CartStore interface {
	CreateCart(userID int) error
	GetMyCartItems(userID int) (*[]*CartItem, error)
	AddItemToCart(productID int, userID int, price float64) (*CartItem, error)
	RemoveItemFromCart(productID int, userID int) error
	GetTotal(userID int) (float64, error)
	GetCartID(userID int) (int, error)
	GetCartItem(userID int, productID int) (*CartItem, error)
	RemoveOneItemFromCart(userID int, productID int) error
	RemoveItemsFromCart(userID int) error
}

type CartService interface {
	CreateCart(userID int) error
	GetMyCartItems(userID int) (*[]*CartItem, error)
	AddItemToCart(productID int, userID int) (*CartItem, error)
	RemoveItemFromCart(productID int, userID int) error
	GetTotal(userID int) (float64, error)
	RemoveEntireItemFromCart(productID int, userID int) error
	RemoveItemsFromCart(userID int) error
}

type Cart struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CartItem struct {
	CartID        int       `json:"cartId"`
	ProductID     int       `json:"productId"`
	ProductTitle  string    `json:"productTitle"`
	Quantity      int       `json:"quantity"`
	ProductImage  string    `json:"productImage"`
	PriceAtAdding float64   `json:"priceAtAdding"`
	AddedAt       time.Time `json:"addedAt"`
}
