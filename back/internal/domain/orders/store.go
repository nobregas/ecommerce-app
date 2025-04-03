package orders

import (
	"database/sql"
	"fmt"

	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) CreateOrder(userID int, totalAmount float64, paymentMethod types.PaymentMethod, paymentID string) (*types.OrderHistory, error) {
	query := `
		INSERT INTO order_history (userId, totalAmount, status, paymentMethod, paymentId, createdAt, updatedAt)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := s.db.Exec(query, userID, totalAmount, types.OrderPending, paymentMethod, paymentID)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting order ID: %w", err)
	}

	order, err := s.GetOrderByID(int(orderID))
	if err != nil {
		return nil, fmt.Errorf("error retrieving created order: %w", err)
	}

	return order, nil
}

func (s *Store) AddOrderItems(orderID int, items []*types.OrderItem) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO order_items (orderId, productId, quantity, price)
		VALUES (?, ?, ?, ?)
	`

	for _, item := range items {
		_, err = tx.Exec(query, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return fmt.Errorf("error adding order item: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (s *Store) GetOrdersByUserID(userID int) ([]*types.OrderHistory, error) {
	query := `
		SELECT id, userId, totalAmount, status, paymentMethod, paymentId, createdAt, updatedAt
		FROM order_history
		WHERE userId = ?
		ORDER BY createdAt DESC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}
	defer rows.Close()

	var orders []*types.OrderHistory
	for rows.Next() {
		order := &types.OrderHistory{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.Status,
			&order.PaymentMethod,
			&order.PaymentID,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

func (s *Store) GetOrderByID(orderID int) (*types.OrderHistory, error) {
	query := `
		SELECT id, userId, totalAmount, status, paymentMethod, paymentId, createdAt, updatedAt
		FROM order_history
		WHERE id = ?
	`
	order := &types.OrderHistory{}
	err := s.db.QueryRow(query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Status,
		&order.PaymentMethod,
		&order.PaymentID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewEntityNotFound("order", orderID)
		}
		return nil, fmt.Errorf("error fetching order: %w", err)
	}

	return order, nil
}

func (s *Store) GetOrderItems(orderID int) ([]*types.OrderItem, error) {
	query := `
		SELECT orderId, productId, quantity, price
		FROM order_items
		WHERE orderId = ?
	`
	rows, err := s.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("error fetching order items: %w", err)
	}
	defer rows.Close()

	var items []*types.OrderItem
	for rows.Next() {
		item := &types.OrderItem{}
		err := rows.Scan(
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order item: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order items: %w", err)
	}

	return items, nil
}

func (s *Store) UpdateOrderStatus(orderID int, status types.OrderStatus) error {
	query := `
		UPDATE order_history
		SET status = ?, updatedAt = NOW()
		WHERE id = ?
	`
	result, err := s.db.Exec(query, status, orderID)
	if err != nil {
		return fmt.Errorf("error updating order status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return apperrors.NewEntityNotFound("order", orderID)
	}

	return nil
}

func (s *Store) GetOrdersWithItems(userID int) ([]*types.OrderWithItems, error) {
	orders, err := s.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	var ordersWithItems []*types.OrderWithItems
	for _, order := range orders {
		items, err := s.GetOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		orderWithItems := &types.OrderWithItems{
			Order: *order,
			Items: items,
		}
		ordersWithItems = append(ordersWithItems, orderWithItems)
	}

	return ordersWithItems, nil
}

func (s *Store) GetOrderWithItems(orderID int) (*types.OrderWithItems, error) {
	order, err := s.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	items, err := s.GetOrderItems(orderID)
	if err != nil {
		return nil, err
	}

	orderWithItems := &types.OrderWithItems{
		Order: *order,
		Items: items,
	}

	return orderWithItems, nil
}
