package discount

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nobregas/ecommerce-mobile-back/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateDiscount(payload *types.CreateProductDiscountPayload) (*types.ProductDiscount, error) {
	query := `
		INSERT INTO product_discounts
			(productId, discountPercent, startDate, endDate)
		Values (?, ?, ?, ?)
	`

	res, err := s.db.ExecContext(
		context.Background(),
		query,
		payload.ProductID,
		payload.DiscountPercent,
		payload.StartDate,
		payload.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create discount: %w", err)
	}

	id, _ := res.LastInsertId()
	discount, err := s.GetDiscoutsByID(int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get discount: %w", err)
	}

	return discount, nil
}

func (s *Store) UpdateDiscount(discountID int, payload *types.UpdateProductDiscountPayload) (*types.ProductDiscount, error) {
	query := `
		UPDATE product_discounts SET
			discountPercent = ?,
			startDate = ?,
			endDate = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		payload.DiscountPercent,
		payload.StartDate,
		payload.EndDate,
		discountID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update discount: %w", err)
	}

	discount, err := s.GetDiscoutsByID(discountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get discount: %w", err)
	}

	return discount, nil
}

func (s *Store) DeleteDiscount(discountID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// deletes discount
	result, err := tx.Exec(`DELETE FROM product_discounts WHERE id = ?`, discountID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// verify if discount exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // discount not found
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil

}

func (s *Store) GetDiscoutsByID(discoutId int) (*types.ProductDiscount, error) {
	query := `
		SELECT id, productId, discountPercent, startDate, endDate, createdAt
		FROM product_discounts
		WHERE id = ?
	`

	row := s.db.QueryRowContext(context.Background(), query, discoutId)

	return scanRowIntoDiscount(row)
}

func (s *Store) GetDiscountsByProduct(productID int) ([]*types.ProductDiscount, error) {
	query := `
		SELECT id, productId, discountPercent, startDate, endDate, createdAt
		FROM product_discounts
		WHERE productId = ?
		ORDER BY startDate DESC
	`

	rows, err := s.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query discounts: %w", err)
	}
	defer rows.Close()

	discounts, err := scanRowsIntoDiscount(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan discounts: %w", err)
	}

	return discounts, nil
}

func (s *Store) GetActiveDiscounts(productID int) ([]*types.ProductDiscount, error) {
	query := `
		SELECT id, productId, discountPercent, startDate, endDate, createdAt
		FROM product_discounts
		WHERE productId = ?
		AND startDate <= NOW()
		AND endDate >= NOW()
	`

	rows, err := s.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query discounts: %w", err)
	}
	defer rows.Close()

	discounts, err := scanRowsIntoDiscount(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan discounts: %w", err)
	}

	return discounts, nil
}

func (s *Store) GetDiscountsByDateRange(productID int, start time.Time, end time.Time) ([]*types.ProductDiscount, error) {
	query := `
		SELECT id, productId, discountPercent, startDate, endDate, createdAt
		FROM product_discounts
		WHERE productId = ?
		AND (
			(startDate BETWEEN ? AND ?) OR
			(endDate BETWEEN ? AND ?) OR
			(startDate <= ? AND endDate >= ?)
		)
	`

	rows, err := s.db.Query(query,
		productID,
		start, end,
		start, end,
		start, end,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query discounts: %w", err)
	}
	defer rows.Close()

	discounts, err := scanRowsIntoDiscount(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan discounts: %w", err)
	}

	return discounts, nil
}

func scanRowsIntoDiscount(rows *sql.Rows) ([]*types.ProductDiscount, error) {
	var discounts []*types.ProductDiscount

	for rows.Next() {
		var discount types.ProductDiscount
		err := rows.Scan(
			&discount.ID,
			&discount.ProductID,
			&discount.DiscountPercent,
			&discount.StartDate,
			&discount.EndDate,
			&discount.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan discount: %w", err)
		}
		discounts = append(discounts, &discount)
	}

	return discounts, nil
}

func scanRowIntoDiscount(row *sql.Row) (*types.ProductDiscount, error) {
	discount := new(types.ProductDiscount)

	err := row.Scan(
		&discount.ID,
		&discount.ProductID,
		&discount.DiscountPercent,
		&discount.StartDate,
		&discount.EndDate,
		&discount.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan discount: %w", err)
	}

	return discount, nil
}
