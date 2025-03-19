package user

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

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (fullName, email, cpf, password) VALUES (?, ?, ?, ?)",
		user.FullName,
		user.Email,
		user.Cpf,
		user.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Cpf,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
