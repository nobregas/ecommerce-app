package favorite

import (
	"database/sql"
	"fmt"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AddFavorite(userID int, productID int) (*types.UserFavorite, error) {
	query := `
		INSERT INTO user_favorites (user_id, product_id)
		VALUES (?, ?)
	`
	_, err := s.db.Exec(query, userID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to add favorite: %w", err)
	}

	createdFavorite, err := s.GetFavorite(userID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created favorite: %w", err)
	}

	return createdFavorite, nil
}

func (s *Store) RemoveFavorite(userID int, productID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// deletes favorite
	result, err := tx.Exec(`DELETE FROM user_favorites WHERE userId = ? AND productId = ?`, userID, productID)
	if err != nil {
		return fmt.Errorf("failed to delete favorite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected in favorite delete: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // rating not found
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction in favorite delete: %w", err)
	}

	return nil
}

func (s *Store) GetUserFavorite(userID int) (*[]*types.UserFavorite, error) {
	query := `SELECT userId, productId, addedAt
		FROM user_favorites
		WHERE userId = ?
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query favorites: %w", err)
	}
	defer rows.Close()

	favorites, err := scanRows(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan favorites: %w", err)
	}
	return favorites, err

}

func (s *Store) GetFavorite(userID int, productID int) (*types.UserFavorite, error) {
	query := `
		SELECT userId, productId, addedAt
		FROM user_favorites
		WHERE user_id = ? AND productId = ?
	`

	row := s.db.QueryRow(query, userID, productID)
	favorite, err := scanRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite: %w", err)
	}
	return favorite, nil
}

func scanRow(row *sql.Row) (*types.UserFavorite, error) {
	userFavorite := new(types.UserFavorite)
	err := row.Scan(
		&userFavorite.UserID,
		&userFavorite.ProductId,
		&userFavorite.AddedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}
	return userFavorite, nil
}

func scanRows(rows *sql.Rows) (*[]*types.UserFavorite, error) {
	var favorites []*types.UserFavorite
	for rows.Next() {
		var r types.UserFavorite
		err := rows.Scan(
			&r.UserID,
			&r.ProductId,
			&r.AddedAt,
		)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, &r)
	}
	return &favorites, nil
}
