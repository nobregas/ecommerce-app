package cart

import (
	"fmt"

	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
)

type Service struct {
	cartStore     types.CartStore
	productStore  types.ProductStore
	discountStore types.ProductDiscountStore
}

func NewService(
	cartStore types.CartStore,
	productStore types.ProductStore,
	discountStore types.ProductDiscountStore,
) *Service {
	return &Service{
		cartStore:     cartStore,
		productStore:  productStore,
		discountStore: discountStore,
	}
}

func (s *Service) CreateCart(userID int) error {
	return s.cartStore.CreateCart(userID)
}

func (s *Service) GetMyCartItems(userID int) (*[]*types.CartItem, error) {
	items, err := s.cartStore.GetMyCartItems(userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) AddItemToCart(productID int, userID int) (*types.CartItem, error) {
	fmt.Printf("[CART SERVICE] Starting to add product %d to cart for user %d\n", productID, userID)

	product, err := s.productStore.GetProductByID(productID)
	if err != nil {
		fmt.Printf("[CART SERVICE] ERROR getting product %d: %v\n", productID, err)
		return nil, err
	}

	if product == nil {
		fmt.Printf("[CART SERVICE] Product %d not found\n", productID)
		return nil, apperrors.NewEntityNotFound("product", productID)
	}

	fmt.Printf("[CART SERVICE] Product %d found. Stock: %d\n", productID, product.Inventory.StockQuantity)
	if product.Inventory.StockQuantity <= 0 {
		fmt.Printf("[CART SERVICE] Product %d out of stock\n", productID)
		return nil, apperrors.NewValidationError("product", "product out of stock")
	}

	fmt.Printf("[CART SERVICE] Checking discounts for product %d\n", productID)
	discounts, err := s.discountStore.GetActiveDiscounts(productID)
	if err != nil {
		fmt.Printf("[CART SERVICE] ERROR getting discounts for product %d: %v\n", productID, err)
		return nil, fmt.Errorf("error getting discounts: %w", err)
	}

	finalPrice := product.BasePrice
	if len(discounts) > 0 {
		highestDiscount := 0.0
		for _, discount := range discounts {
			if discount.DiscountPercent > highestDiscount {
				highestDiscount = discount.DiscountPercent
			}
		}
		fmt.Printf("[CART SERVICE] Applying discount %.2f%% to product %d. Original: %.2f, Final: %.2f\n",
			highestDiscount, productID, product.BasePrice, finalPrice)

		finalPrice = product.BasePrice * (1 - highestDiscount/100)
	}

	fmt.Printf("[CART SERVICE] Sending to store: product %d, user %d, price %.2f\n", productID, userID, finalPrice)
	item, err := s.cartStore.AddItemToCart(productID, userID, finalPrice)
	if err != nil {
		fmt.Printf("[CART SERVICE] ERROR adding item to cart: %v\n", err)
		return nil, err
	}

	fmt.Printf("[CART SERVICE] Successfully added. New quantity: %d\n", item.Quantity)
	return item, nil
}

func (s *Service) RemoveItemFromCart(productID int, userID int) error {

	item, err := s.cartStore.GetCartItem(userID, productID)
	if err != nil {
		fmt.Printf("[CART SERVICE] ERROR getting item at remove item from cart %d: %v\n", productID, err)
		return err
	}

	if item.Quantity > 1 {
		err := s.cartStore.RemoveOneItemFromCart(userID, productID)
		if err != nil {
			fmt.Printf("[CART SERVICE] ERROR removing one item from cart %d: %v\n", productID, err)
			return err
		}
	} else {
		err := s.cartStore.RemoveItemFromCart(userID, productID)
		if err != nil {
			fmt.Printf("[CART SERVICE] ERROR removing item from cart %d: %v\n", productID, err)
			return err
		}
	}

	return nil
}

func (s *Service) RemoveEntireItemFromCart(productID int, userID int) error {
	fmt.Printf("[CART SERVICE] removing product %d from user %d\n", productID, userID)
	return s.cartStore.RemoveItemFromCart(productID, userID)
}

func (s *Service) GetTotal(userID int) (float64, error) {
	total, err := s.cartStore.GetTotal(userID)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Service) RemoveItemsFromCart(userID int) error {
	fmt.Printf("[CART SERVICE] Removing all items from cart for user %d\n", userID)
	return s.cartStore.RemoveItemsFromCart(userID)
}
