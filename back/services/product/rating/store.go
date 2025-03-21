package rating

import (
	"database/sql"
	"fmt"

	"github.com/nobregas/ecommerce-mobile-back/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateRating(payload *types.CreateProductRatingPayload, userID int) (*types.ProductRating, error) {
	query := `
		INSERT INTO product_ratings
			(userId, productId, rating, comment)
		VALUES (?, ?, ?, ?)
	`

	res, err := s.db.Exec(query,
		userID,
		payload.ProductID,
		payload.Rating,
		payload.Comment)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	rating, err := s.GetRating(int(id))
	if err != nil {
		return nil, err
	}

	return rating, nil
}

func (s *Store) GetRatingsByProduct(productID int) ([]*types.ProductRating, error) {
	query := `
		SELECT id, userId, productId, rating, comment, createdAt
		FROM product_ratings
		WHERE productId = ?
		ORDER BY createdAt DESC
	`

	rows, err := s.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRatings(rows)
}

func (s *Store) GetRatingsByUser(int) ([]*types.ProductRating, error) {
	query := `
		SELECT id, userId, productId, rating, comment, createdAt
		FROM product_ratings
		WHERE userId = ?
		ORDER BY createdAt DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRatings(rows)
}

func (s *Store) GetRating(ratingID int) (*types.ProductRating, error) {
	query := `
		SELECT id, userId, productId, rating, comment, createdAt
		FROM product_ratings
		WHERE id = ?
	`

	row := s.db.QueryRow(query, ratingID)
	rating, err := scanRating(row)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

func (s *Store) UpdateRating(ratingID int, payload *types.UpdateProductRatingPayload) (*types.ProductRating, error) {
	query := `
		UPDATE product_ratings SET
			rating = ?,
			comment = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(query, payload.Rating, payload.Comment, ratingID)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while updating a rating: %w", err)
	}

	rating, err := s.GetRating(ratingID)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while getting rating in update: %w", err)
	}

	return rating, nil
}

func (s *Store) DeleteRating(ratingID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// deletes rating
	result, err := tx.Exec(`DELETE FROM product_ratings WHERE id = ?`, ratingID)
	if err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}

	// verify if rating exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected in rating delete: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // rating not found
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction in rating delete: %w", err)
	}

	return nil
}

func (s *Store) GetAverageRating(productID int) (float64, error) {
	query := `
        SELECT COALESCE(AVG(rating), 0) 
        FROM product_ratings 
        WHERE productId = ?`

	avg := new(float64)

	err := s.db.QueryRow(query, productID).Scan(avg)
	if err != nil {
		return -1, fmt.Errorf("failed to get average rating: %w", err)
	}

	return *avg, nil
}

func scanRating(row *sql.Row) (*types.ProductRating, error) {
	r := new(types.ProductRating)
	err := row.Scan(
		&r.ID,
		&r.UserID,
		&r.ProductID,
		&r.Rating,
		&r.Comment,
		&r.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rating: %w", err)
	}

	return r, nil
}

func scanRatings(rows *sql.Rows) ([]*types.ProductRating, error) {
	var ratings []*types.ProductRating
	for rows.Next() {
		var r types.ProductRating
		err := rows.Scan(
			&r.ID,
			&r.UserID,
			&r.ProductID,
			&r.Rating,
			&r.Comment,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, &r)
	}
	return ratings, nil
}
