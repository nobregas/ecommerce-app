package product

import (
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
)

type ProductService struct {
	productStore  types.ProductStore
	userStore     types.UserStore
	discountStore types.ProductDiscountStore
	ratingStore   types.ProductRatingStore
}

func NewProductService(
	productStore types.ProductStore,
	userStore types.UserStore,
	discountStore types.ProductDiscountStore,
	ratingStore types.ProductRatingStore) *ProductService {

	return &ProductService{
		productStore:  productStore,
		userStore:     userStore,
		discountStore: discountStore,
		ratingStore:   ratingStore,
	}
}

func (p *ProductService) GetProductDetails() *types.ProductDetails {
	return nil
}

func (p *ProductService) GetSimpleProducts() *[]types.SimpleProductObject {

	return nil
}

func (p *ProductService) GetProducts() []*types.Product {
	products, err := p.productStore.GetProducts()
	if err != nil {
		panic(err)
		return nil
	}

	return products
}

func (p *ProductService) GetProductByID(productID int) *types.Product {
	product, err := p.productStore.GetProductByID(productID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("Product", productID))
	}

	return product
}

func (p *ProductService) GetProductsByCategoryID(categoryID int) []*types.Product {
	products, err := p.productStore.GetProductsByCategory(categoryID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("Category", categoryID))
		return nil
	}

	return products
}

func (p *ProductService) CreateProductWithImages(payload types.CreateProductWithImagesPayload) *types.Product {
	if err := utils.Validate.Struct(payload); err != nil {
		panic(apperrors.NewValidationError("invalid payload", err.Error()))
		return nil
	}

	createdProduct, err := p.productStore.CreateProductWithImages(payload)
	if err != nil {
		panic(err)
		return nil
	}

	return createdProduct
}

func (p *ProductService) UpdateProductById(productID int, payload types.UpdateProductPayload) *types.Product {
	_, err := p.productStore.GetProductByID(productID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("Product", productID))
		return nil
	}

	if err := utils.Validate.Struct(payload); err != nil {
		panic(apperrors.NewValidationError("invalid payload", err.Error()))
		return nil
	}

	if err := p.productStore.UpdateProduct(productID, payload); err != nil {
		panic(err)
		return nil
	}

	updatedProduct, err := p.productStore.GetProductByID(productID)
	if err != nil {
		panic(err)
		return nil
	}

	return updatedProduct
}

func (p *ProductService) DeleteProduct(productID int) {
	_, err := p.productStore.GetProductByID(productID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("Product", productID))
	}

	if err := p.productStore.DeleteProduct(productID); err != nil {
		panic(err)
	}
}
