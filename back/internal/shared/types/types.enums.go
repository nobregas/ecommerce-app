package types

import (
	"database/sql/driver"
	"fmt"
)

// UserRole
type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

func (r UserRole) Valid() error {
	switch r {
	case RoleUser, RoleAdmin:
		return nil
	default:
		return fmt.Errorf("invalid user role: %s", r)
	}
}

func (r UserRole) Value() (driver.Value, error) {
	if err := r.Valid(); err != nil {
		return nil, err
	}
	return string(r), nil
}

func (r *UserRole) Scan(value interface{}) error {
	if value == nil {
		*r = RoleUser // default
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for UserRole")
	}

	*r = UserRole(str)
	return r.Valid()
}

// OrderStatus
type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderCompleted OrderStatus = "COMPLETED"
	OrderCancelled OrderStatus = "CANCELLED"
	OrderShipped   OrderStatus = "SHIPPED"
)

func (s OrderStatus) Valid() error {
	switch s {
	case OrderPending, OrderCompleted, OrderCancelled, OrderShipped:
		return nil
	default:
		return fmt.Errorf("invalid order status: %s", s)
	}
}

// PaymentMethod
type PaymentMethod string

const (
	PaymentCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentDebitCard    PaymentMethod = "DEBIT_CARD"
	PaymentPix          PaymentMethod = "PIX"
	PaymentBankTransfer PaymentMethod = "BANK_TRANSFER"
)

func (p PaymentMethod) Valid() error {
	switch p {
	case PaymentCreditCard, PaymentDebitCard, PaymentPix, PaymentBankTransfer:
		return nil
	default:
		return fmt.Errorf("invalid payment method: %s", p)
	}
}
