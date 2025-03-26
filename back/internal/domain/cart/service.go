package cart

import "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"

type Service struct {
	store types.CartStore
}

func NewService(store types.CartStore) *Service {
	return &Service{store: store}
}
