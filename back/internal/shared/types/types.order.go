package types

type OrderStore interface {
	CreateOrder(userID int, totalAmount float64, paymentMethod PaymentMethod, paymentID string) (*OrderHistory, error)
	AddOrderItems(orderID int, items []*OrderItem) error
	GetOrdersByUserID(userID int) ([]*OrderHistory, error)
	GetOrderByID(orderID int) (*OrderHistory, error)
	GetOrderItems(orderID int) ([]*OrderItem, error)
	UpdateOrderStatus(orderID int, status OrderStatus) error
	GetOrdersWithItems(userID int) ([]*OrderWithItems, error)
	GetOrderWithItems(orderID int) (*OrderWithItems, error)
}

type OrderService interface {
	CreateOrderFromCart(userID int, paymentMethod PaymentMethod, paymentID string) (*OrderHistory, error)
	GetOrdersByUserID(userID int) ([]*OrderHistory, error)
	GetOrderByID(orderID int) (*OrderHistory, error)
	GetOrderWithItems(orderID int) (*OrderWithItems, error)
	GetOrdersWithItems(userID int) ([]*OrderWithItems, error)
	UpdateOrderStatus(orderID int, status OrderStatus) error
}

type OrderWithItems struct {
	Order OrderHistory `json:"order"`
	Items []*OrderItem `json:"items"`
}

type CreateOrderPayload struct {
	PaymentMethod PaymentMethod `json:"paymentMethod" validate:"required"`
	PaymentID     string        `json:"paymentId"`
}

type UpdateOrderStatusPayload struct {
	Status OrderStatus `json:"status" validate:"required"`
}
