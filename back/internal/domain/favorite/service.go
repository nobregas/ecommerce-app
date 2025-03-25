package favorite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
)

type Service struct {
	favoriteStore types.UserFavoriteStore
	userStore     types.UserStore
	productStore  types.ProductStore
}

func NewService(
	favoriteStore types.UserFavoriteStore,
	userStore types.UserStore,
	productStore types.ProductStore,
) *Service {
	return &Service{
		favoriteStore: favoriteStore,
		userStore:     userStore,
		productStore:  productStore,
	}
}

func (s *Service) AddFavorite(userID int, productID int) *types.UserFavorite {
	_, err := s.userStore.GetUserByID(userID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("user", userID))
	}

	_, err = s.productStore.GetProductByID(productID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("product", productID))
	}

	existing, err := s.favoriteStore.GetFavorite(userID, productID)
	if existing != nil {
		panic(apperrors.NewConflictError("product", "product already favorited"))
	}

	favorite, err := s.favoriteStore.AddFavorite(userID, productID)
	if err != nil {
		panic(fmt.Errorf("failed to add favorite: %w", err))
	}

	return favorite
}

func (s *Service) RemoveFavorite(userID int, productID int) {
	_, err := s.userStore.GetUserByID(userID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("user", userID))
	}

	err = s.favoriteStore.RemoveFavorite(userID, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(apperrors.NewEntityNotFound("favorite", fmt.Sprintf("user %d product %d", userID, productID)))
		}
		panic(fmt.Errorf("failed to remove favorite: %w", err))
	}
}

func (s *Service) GetUserFavorite(userID int) *[]*types.UserFavorite {
	_, err := s.userStore.GetUserByID(userID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("user", userID))
	}

	favorites, err := s.favoriteStore.GetUserFavorite(userID)
	if err != nil {
		panic(fmt.Errorf("failed to get favorites: %w", err))
	}

	return favorites
}

func (s *Service) IsFavorited(userID int, productID int) bool {
	isFav, err := s.favoriteStore.IsFavorited(userID, productID)
	if err != nil {
		panic(fmt.Errorf("failed to check favorite status: %w", err))
	}
	return isFav
}
