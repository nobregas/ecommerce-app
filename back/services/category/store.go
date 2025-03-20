package category

import (
	"database/sql"

	"github.com/nobregas/ecommerce-mobile-back/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetCategories() (*[]types.Category, error) {
	return nil, nil
}

func (s *Store) GetCategoryByID(categoryID int) (*types.Category, error) {
	return nil, nil
}

func (s *Store) CreateCategory(payload types.CreateCategoryPayload) (*types.Category, error) {
	return nil, nil
}

func (s *Store) UpdateCategory(categoryID int, payload types.UpdateCategoryPayload) (*types.Category, error) {
	return nil, nil
}

func (s *Store) DeleteCategory(categoryID int) error {
	return nil
}
