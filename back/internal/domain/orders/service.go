package orders

import (
	"fmt"

	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
)

type Service struct {
	orderStore   types.OrderStore
	cartStore    types.CartStore
	productStore types.ProductStore
}

func NewService(
	orderStore types.OrderStore,
	cartStore types.CartStore,
	productStore types.ProductStore,
) *Service {
	return &Service{
		orderStore:   orderStore,
		cartStore:    cartStore,
		productStore: productStore,
	}
}

func (s *Service) CreateOrderFromCart(userID int, paymentMethod types.PaymentMethod, paymentID string) (*types.OrderHistory, error) {
	fmt.Printf("[ORDER SERVICE] Creating order for user %d with payment method %s\n", userID, paymentMethod)

	if err := paymentMethod.Valid(); err != nil {
		return nil, apperrors.NewValidationError("paymentMethod", err.Error())
	}

	cartItems, err := s.cartStore.GetMyCartItems(userID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error getting cart items: %v\n", err)
		return nil, fmt.Errorf("error getting cart items: %w", err)
	}

	if cartItems == nil || len(*cartItems) == 0 {
		return nil, apperrors.NewValidationError("cart", "cart is empty")
	}

	total, err := s.cartStore.GetTotal(userID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error calculating cart total: %v\n", err)
		return nil, fmt.Errorf("error calculating cart total: %w", err)
	}

	order, err := s.orderStore.CreateOrder(userID, total, paymentMethod, paymentID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error creating order: %v\n", err)
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	var orderItems []*types.OrderItem
	for _, cartItem := range *cartItems {
		orderItem := &types.OrderItem{
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.PriceAtAdding,
		}
		orderItems = append(orderItems, orderItem)

		err = s.productStore.UpdateStock(cartItem.ProductID, -cartItem.Quantity)
		if err != nil {
			fmt.Printf("[ORDER SERVICE] Error updating stock for product %d: %v\n", cartItem.ProductID, err)
			return nil, fmt.Errorf("error updating stock: %w", err)
		}
	}

	err = s.orderStore.AddOrderItems(order.ID, orderItems)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error adding order items: %v\n", err)
		return nil, fmt.Errorf("error adding order items: %w", err)
	}

	err = s.cartStore.RemoveItemsFromCart(userID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error clearing cart: %v\n", err)
		return nil, fmt.Errorf("error clearing cart: %w", err)
	}

	fmt.Printf("[ORDER SERVICE] Order created successfully with ID %d\n", order.ID)
	return order, nil
}

func (s *Service) GetOrdersByUserID(userID int) ([]*types.OrderHistory, error) {
	fmt.Printf("[ORDER SERVICE] Getting orders for user %d\n", userID)

	orders, err := s.orderStore.GetOrdersByUserID(userID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error getting orders: %v\n", err)
		return nil, fmt.Errorf("error getting orders: %w", err)
	}

	return orders, nil
}

func (s *Service) GetOrderByID(orderID int) (*types.OrderHistory, error) {
	fmt.Printf("[ORDER SERVICE] Getting order with ID %d\n", orderID)

	order, err := s.orderStore.GetOrderByID(orderID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error getting order: %v\n", err)
		return nil, err
	}

	return order, nil
}

func (s *Service) GetOrderWithItems(orderID int) (*types.OrderWithItems, error) {
	fmt.Printf("[ORDER SERVICE] Getting order with items for order ID %d\n", orderID)

	orderWithItems, err := s.orderStore.GetOrderWithItems(orderID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error getting order with items: %v\n", err)
		return nil, err
	}

	return orderWithItems, nil
}

func (s *Service) GetOrdersWithItems(userID int) ([]*types.OrderWithItems, error) {
	fmt.Printf("[ORDER SERVICE] Getting orders with items for user ID %d\n", userID)

	ordersWithItems, err := s.orderStore.GetOrdersWithItems(userID)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error getting orders with items: %v\n", err)
		return nil, err
	}

	return ordersWithItems, nil
}

func (s *Service) UpdateOrderStatus(orderID int, status types.OrderStatus) error {
	fmt.Printf("[ORDER SERVICE] Updating order %d status to %s\n", orderID, status)

	if err := status.Valid(); err != nil {
		return apperrors.NewValidationError("status", err.Error())
	}

	err := s.orderStore.UpdateOrderStatus(orderID, status)
	if err != nil {
		fmt.Printf("[ORDER SERVICE] Error updating order status: %v\n", err)
		return err
	}

	return nil
}
