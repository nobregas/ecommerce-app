package types

import "time"

type ProductStore interface {
	GetProducts() ([]*Product, error)
	CreateProduct(CreateProductPayload) error
	GetProductByID(productID int) (*Product, error)
	CreateProductWithImages(CreateProductWithImagesPayload) (*Product, error)
	UpdateStock(productID int, quantityChange int) error
	GetInventory(productID int) (*Inventory, error)
	GetImagesForProducts(productIDs []int) (map[int][]ProductImage, error)
	UpdateProduct(productID int, payload UpdateProductPayload) error
	DeleteProduct(productID int) error
	GetProductsByCategory(categoryID int) ([]*Product, error)
	GetProductDetails(userID int, productID int) (*ProductDetails, error)
	GetSimpleProductDetails(userID int) (*[]*SimpleProductObject, error)
}

type ProductService interface {
	GetProductDetails(userID int, productID int) *ProductDetails
	GetSimpleProducts(userID int) *[]*SimpleProductObject
	GetProducts() []*Product
	GetProductByID(productID int) *Product
	GetProductsByCategoryID(categoryID int) []*Product
	CreateProductWithImages(payload CreateProductWithImagesPayload) *Product
	UpdateProductById(productID int, payload UpdateProductPayload) *Product
	DeleteProduct(productID int)
}

type Product struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	BasePrice   float64        `json:"basePrice"`
	Inventory   Inventory      `json:"inventory"`
	Images      []ProductImage `json:"images"`
	Categories  []Category     `json:"categories"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type Inventory struct {
	ProductID     int `json:"productId"`
	StockQuantity int `json:"stockQuantity"`
	Version       int `json:"-"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ProductID int    `json:"productId"`
	ImageUrl  string `json:"imageUrl"`
	SortOrder int    `json:"sortOrder"`
}

type CreateProductPayload struct {
	Title         string  `json:"title" validate:"required,min=3,max=100"`
	Description   string  `json:"description" validate:"max=1000"`
	BasePrice     float64 `json:"basePrice" validate:"required,gt=0"`
	StockQuantity int     `json:"stockQuantity" validate:"required,min=0"`
	CategoryIDs   []int   `json:"categoryIds" validate:"required,min=1"`
}

type CreateProductWithImagesPayload struct {
	Title         string         `json:"title" validate:"required,min=3,max=100"`
	Description   string         `json:"description" validate:"max=1000"`
	BasePrice     float64        `json:"basePrice" validate:"required,gt=0"`
	StockQuantity int            `json:"stockQuantity" validate:"required,min=0"`
	Images        []ImagePayload `json:"images" validate:"required,min=1,dive"`
	CategoryIDs   []int          `json:"categoryIds" validate:"required,min=1"`
}

type UpdateProductPayload struct {
	Title       *string              `json:"title,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string              `json:"description,omitempty" validate:"omitempty,max=1000"`
	BasePrice   *float64             `json:"basePrice,omitempty" validate:"omitempty,gt=0"`
	Images      []ImageUpdatePayload `json:"images,omitempty" validate:"omitempty,dive"`
}

type ImagePayload struct {
	ImageUrl  string `json:"imageUrl" validate:"required"`
	SortOrder int    `json:"sortOrder" validate:"min=0"`
}

type ImageUpdatePayload struct {
	ID        *int   `json:"id,omitempty"`
	ImageUrl  string `json:"imageUrl" validate:"required"`
	SortOrder int    `json:"sortOrder" validate:"min=0"`
	Deleted   bool   `json:"deleted,omitempty"`
}

type ProductDetails struct {
	ID                 *int           `json:"id"`
	Title              string         `json:"title"`
	Price              float64        `json:"price"`
	BasePrice          float64        `json:"basePrice"`
	DiscountPercentage float64        `json:"discountPercentage"`
	Description        string         `json:"description"`
	IsFavorite         bool           `json:"isFavorite"`
	AverageRating      float64        `json:"averageRating"`
	Images             []ProductImage `json:"images"`
}

type SimpleProductObject struct {
	ID            int          `json:"id"`
	Title         string       `json:"title"`
	Price         float64      `json:"price"`
	BasePrice     float64      `json:"basePrice"`
	AverageRating float64      `json:"averageRating"`
	Image         ProductImage `json:"image"`
	IsFavorite    bool         `json:"isFavorite"`
}
