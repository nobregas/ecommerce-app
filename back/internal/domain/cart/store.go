package cart

import (
	"database/sql"
	"errors"
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

func (s *Store) CreateCart(userID int) error {
	query := `
		INSERT INTO carts (userId, createdAt, updatedAt)
		VALUES (?, NOW(), NOW())
	`
	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error creating cart: %w", err)
	}
	return nil
}

func (s *Store) GetMyCartItems(userID int) (*[]*types.CartItem, error) {
	cartID, err := s.GetCartID(userID)
	if err != nil {
		return &[]*types.CartItem{}, fmt.Errorf("error getting cart ID: %w", err)
	}

	query := `
		SELECT 
			c.id as cartId,
			ci.productId,
			ci.quantity,
			ci.priceAtAdding,
			ci.addedAt,
			ci.productImage,
			ci.productTitle
		FROM carts c
		JOIN cart_items ci ON c.id = ci.cartId
		WHERE c.id = ?
	`
	rows, err := s.db.Query(query, cartID)
	if err != nil {
		return &[]*types.CartItem{}, fmt.Errorf("error fetching cart items: %w", err)
	}
	defer rows.Close()

	items, err := scanRows(rows)
	if err != nil {
		return &[]*types.CartItem{}, fmt.Errorf("error scanning cart items: %w", err)
	}

	return items, nil
}

func (s *Store) AddItemToCart(productID int, userID int, price float64) (*types.CartItem, error) {
	cartID, err := s.GetCartID(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting cart ID: %w", err)
	}

	var productTitle string
	var productImage string
	err = s.db.QueryRow(`
		SELECT 
			p.title,
			COALESCE(pi.imageUrl, '') 
		FROM products p
		LEFT JOIN product_images pi 
			ON p.id = pi.productId 
			AND pi.sortOrder = (SELECT MIN(sortOrder) FROM product_images WHERE productId = p.id)
		WHERE p.id = ?
	`, productID).Scan(&productTitle, &productImage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewEntityNotFound("product", productID)
		}
		return nil, fmt.Errorf("error fetching product details: %w", err)
	}

	result, err := s.db.Exec(`
		UPDATE cart_items 
		SET quantity = quantity + 1 
		WHERE cartId = ? AND productId = ?
	`, cartID, productID)
	if err != nil {
		return nil, fmt.Errorf("error updating cart item: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
		row := s.db.QueryRow(`
			SELECT 
				cartId, 
				productId, 
				quantity, 
				priceAtAdding, 
				addedAt, 
				productImage,
				productTitle 
			FROM cart_items 
			WHERE cartId = ? AND productId = ?
		`, cartID, productID)
		return scanRow(row)
	}

	_, err = s.db.Exec(`
		INSERT INTO cart_items 
			(cartId, productId, quantity, priceAtAdding, addedAt, productImage, productTitle)
		VALUES (?, ?, 1, ?, NOW(), ?, ?)
	`, cartID, productID, price, productImage, productTitle)
	if err != nil {
		return nil, fmt.Errorf("error inserting cart item: %w", err)
	}

	newItemRow := s.db.QueryRow(`
		SELECT 
			cartId, 
			productId, 
			quantity, 
			priceAtAdding, 
			addedAt, 
			productImage,
			productTitle 
		FROM cart_items 
		WHERE cartId = ? AND productId = ?
	`, cartID, productID)

	return scanRow(newItemRow)
}

func (s *Store) RemoveItemFromCart(productID int, userID int) error {
	cartID, err := s.GetCartID(userID)
	if err != nil {
		return fmt.Errorf("error getting cart ID: %w", err)
	}

	fmt.Printf("[CART STORE] removing product %d from user %d and cart %d\n", productID, userID, cartID)

	result, err := s.db.Exec(`
		DELETE FROM cart_items 
		WHERE cartId = ? AND productId = ?
	`, cartID, productID)
	if err != nil {
		return fmt.Errorf("error removing item from cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking removed item: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found in cart")
	}

	return nil
}

func (s *Store) GetTotal(userID int) (float64, error) {
	cartID, err := s.GetCartID(userID)
	if err != nil {
		return 0, fmt.Errorf("error getting cart ID: %w", err)
	}

	var total float64
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(priceAtAdding * quantity), 0)
		FROM cart_items 
		WHERE cartId = ?
	`, cartID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("error calculating cart total: %w", err)
	}
	return total, nil
}

func (s *Store) GetCartID(userID int) (int, error) {
	var cartID int
	err := s.db.QueryRow(`
		SELECT id 
		FROM carts 
		WHERE userId = ?
	`, userID).Scan(&cartID)

	if errors.Is(err, sql.ErrNoRows) {
		res, err := s.db.Exec(`
			INSERT INTO carts (userId, createdAt, updatedAt) 
			VALUES (?, NOW(), NOW())
		`, userID)
		if err != nil {
			return 0, fmt.Errorf("error creating cart: %w", err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("error getting cart ID: %w", err)
		}
		return int(id), nil
	}

	if err != nil {
		return 0, fmt.Errorf("error getting cart ID: %w", err)
	}

	return cartID, nil
}

func (s *Store) GetCartItem(userID int, productID int) (*types.CartItem, error) {
	cartID, err := s.GetCartID(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting cart ID: %w", err)
	}

	row := s.db.QueryRow(`
		SELECT 
			cartId, 
			productId, 
			quantity, 
			priceAtAdding, 
			addedAt, 
			productImage,
			productTitle 
		FROM cart_items 
		WHERE cartId = ? AND productId = ?
	`, cartID, productID)

	return scanRow(row)
}

func (s *Store) RemoveOneItemFromCart(userID int, productID int) error {
	cartID, err := s.GetCartID(userID)
	if err != nil {
		fmt.Printf("[CART STORE]: ERROR getting cart ID at removeOneItemFromCart: %v", err)
		return err
	}

	query := `UPDATE cart_items 
		SET quantity = quantity - 1 
		WHERE cartId = ? AND productId = ? AND quantity > 1;
	`
	_, err = s.db.Exec(query, cartID, productID)
	if err != nil {
		fmt.Printf("[CART STORE]: ERROR removing one item from cart: %w", err)
		return err
	}

	return nil
}

func (s *Store) RemoveItemsFromCart(userID int) error {
	cartID, err := s.GetCartID(userID)
	if err != nil {
		return fmt.Errorf("error getting cart ID: %w", err)
	}

	_, err = s.db.Exec(`DELETE FROM cart_items WHERE cartId = ?`, cartID)
	if err != nil {
		return fmt.Errorf("error deleting cart items: %w", err)
	}

	return nil
}

func scanRow(row *sql.Row) (*types.CartItem, error) {
	c := new(types.CartItem)
	err := row.Scan(
		&c.CartID,
		&c.ProductID,
		&c.Quantity,
		&c.PriceAtAdding,
		&c.AddedAt,
		&c.ProductImage,
		&c.ProductTitle,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}
	return c, nil
}

func scanRows(rows *sql.Rows) (*[]*types.CartItem, error) {
	carts := make([]*types.CartItem, 0)

	for rows.Next() {
		c := new(types.CartItem)
		err := rows.Scan(
			&c.CartID,
			&c.ProductID,
			&c.Quantity,
			&c.PriceAtAdding,
			&c.AddedAt,
			&c.ProductImage,
			&c.ProductTitle,
		)
		if err != nil {
			return &carts, err
		}
		carts = append(carts, c)
	}

	if err := rows.Err(); err != nil {
		return &carts, err
	}

	return &carts, nil
}
