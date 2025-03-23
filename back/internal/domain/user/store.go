package user

import (
	"database/sql"
	"fmt"
	types "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"log"
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
			log.Fatal(err)
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	query := `
        SELECT id, fullName, email, cpf, password, createdAt, updatedAt, role, profile_img 
        FROM users 
        WHERE id = ?`

	row := s.db.QueryRow(query, id)
	return s.scanRowIntoUser(row)
}

func (s *Store) GetUserByCPF(cpf string) (*types.User, error) {
	query := `
        SELECT id, fullName, email, cpf, password, createdAt, updatedAt, role, profile_img 
        FROM users 
        WHERE cpf = ?`

	row := s.db.QueryRow(query, cpf)
	return s.scanRowIntoUser(row)
}

func (s *Store) CreateUser(user types.User) error {
	query := `
        INSERT INTO users 
            (fullName, email, cpf, password, role, profile_img) 
        VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		user.FullName,
		user.Email,
		user.Cpf,
		user.Password,
		user.Role,
		user.ProfileImg,
	)
	return err
}

func (s *Store) scanRowIntoUser(row *sql.Row) (*types.User, error) {
	user := &types.User{}
	var roleStr string

	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Cpf,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&roleStr,
		&user.ProfileImg,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error scanning row: %w", err)
	}

	user.Role = types.UserRole(roleStr)
	if err := user.Role.Valid(); err != nil {
		return nil, fmt.Errorf("invalid role value '%s': %w", roleStr, err)
	}

	return user, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	var roleStr string

	err := rows.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Cpf,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&roleStr,
		&user.ProfileImg,
	)
	if err != nil {
		return nil, err
	}

	user.Role = types.UserRole(roleStr)
	if err := user.Role.Valid(); err != nil {
		return nil, fmt.Errorf("invalid role value '%s': %w", roleStr, err)
	}
	return user, nil
}
