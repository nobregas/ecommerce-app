package category

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/nobregas/ecommerce-mobile-back/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetCategories() (*[]types.Category, error) {
	rows, err := s.db.Query(`
		SELECT id, name, imageUrl, parentCategoryId 
		FROM categories
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	categories := make([]types.Category, 0)
	for rows.Next() {
		c, err := scanRowsIntoCategory(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, *c)
	}

	return &categories, nil
}

func (s *Store) GetCategoryByID(categoryID int) (*types.Category, error) {
	row := s.db.QueryRow(`
		SELECT id, name, imageUrl, parentCategoryId 
		FROM categories 
		WHERE id = ?
	`, categoryID)

	c, err := scanRowIntoCategory(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return c, nil
}

func (s *Store) CreateCategory(payload types.CreateCategoryPayload) (*types.Category, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// validate if parent category exists
	if payload.ParentCategoryId != nil {
		_, err := s.GetCategoryByID(*payload.ParentCategoryId)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent category: %w", err)
		}
	}

	// create category
	res, err := tx.Exec(`
		INSERT INTO categories (name, imageUrl, parentCategoryId) 
		VALUES (?, ?, ?)
	`, payload.Name, payload.ImageUrl, payload.ParentCategoryId)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// get category id
	categoryID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get category id: %w", err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// return category
	return s.GetCategoryByID(int(categoryID))
}

func (s *Store) UpdateCategory(categoryID int, payload types.UpdateCategoryPayload) (*types.Category, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// update category
	query := "UPDATE categories SET"
	args := make([]interface{}, 0)
	updates := []string{}

	if payload.Name != "" {
		updates = append(updates, " name = ?")
		args = append(args, payload.Name)
	}
	if payload.ImageUrl != "" {
		updates = append(updates, " imageUrl = ?")
		args = append(args, payload.ImageUrl)
	}
	if payload.ParentCategoryId != nil {
		updates = append(updates, " parentCategoryId = ?")
		args = append(args, *payload.ParentCategoryId)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(updates, ",") + " WHERE id = ?"
	args = append(args, categoryID)

	// update category
	_, err = tx.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.GetCategoryByID(categoryID)
}

func (s *Store) DeleteCategory(categoryID int) error {
	// delete category
	result, err := s.db.Exec(`
		DELETE FROM categories 
		WHERE id = ?
	`, categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	// verify if category exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func scanRowsIntoCategory(rows *sql.Rows) (*types.Category, error) {
	category := new(types.Category)

	err := rows.Scan(
		&category.ID,
		&category.Name,
		&category.ImageUrl,
		&category.ParentCategoryId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan category: %w", err)
	}

	return category, nil
}

func scanRowIntoCategory(rows *sql.Row) (*types.Category, error) {
	category := new(types.Category)

	err := rows.Scan(
		&category.ID,
		&category.Name,
		&category.ImageUrl,
		&category.ParentCategoryId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan category: %w", err)
	}

	return category, nil
}
