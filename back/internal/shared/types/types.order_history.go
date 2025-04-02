package types

import "time"

type OrderHistory struct {
	ID            int           `json:"id"`
	UserID        int           `json:"userId"`
	TotalAmount   float64       `json:"totalAmount"`
	Status        OrderStatus   `json:"status"`
	PaymentMethod PaymentMethod `json:"paymentMethod"`
	PaymentID     string        `json:"paymentId"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

type OrderItem struct {
	OrderID   int     `json:"orderId"`
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
