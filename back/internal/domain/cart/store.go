package cart

import (
	"database/sql"
	"fmt"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
)

type Store struct {
	*sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) CreateCart(userID int) error {
	query := `
		INSERT INTO carts (userId)
		VALUES (?)
	`
	_, err := s.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("create cart: %w", err)
	}
	return nil
}

func (s *Store) GetMyCartItems() (*[]types.Cartitem, error) {
	return nil, nil
}
func (s *Store) AddItemToCart(productID int, userID int) (*types.Cartitem, error) {
	return nil, nil
}
func (s *Store) RemoveItemFromCart(productID int, userID int) error {
	return nil
}
func (s *Store) GetTotal(userID int) (float64, error) {
	return 0, nil
}
