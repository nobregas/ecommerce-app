package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	GetUserByCPF(cpf string) (*User, error)
	CreateUser(User) error
}

type User struct {
	ID         int       `json:"id"`
	FullName   string    `json:"fullName"`
	Email      string    `json:"email"`
	Cpf        string    `json:"cpf"`
	Role       UserRole  `json:"role" validate:"required,oneof=USER ADMIN"`
	ProfileImg string    `json:"profileImg"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json: "createdAt"`
	UpdatedAt  time.Time `json: "updatedAt"`
}

type RegisterUserPayload struct {
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Cpf      string `json:"cpf" validate:"required"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *User) Sanitize() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"fullName":   u.FullName,
		"email":      u.Email,
		"createdAt":  u.CreatedAt,
		"role":       u.Role,
		"profileImg": u.ProfileImg,
	}
}
