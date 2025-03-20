package types

type CategoryStore interface {
	GetCategories() (*[]Category, error)
	GetCategoryByID(int) (*Category, error)
	CreateCategory(CreateCategoryPayload) (*Category, error)
	UpdateCategory(int, UpdateCategoryPayload) (*Category, error)
	DeleteCategory(int) error
}

type Category struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ImageUrl         string `json:"imageUrl"`
	ParentCategoryId *int   `json:"parentCategoryId"`
}

type CreateCategoryPayload struct {
	Name             string `json:"name" validate:"required,min=3,max=100"`
	ImageUrl         string `json:"imageUrl" validate:"required"`
	ParentCategoryId *int   `json:"parentCategoryId"`
}

type UpdateCategoryPayload struct {
	Name             string `json:"name" validate:"required,min=3,max=100"`
	ImageUrl         string `json:"imageUrl" validate:"required"`
	ParentCategoryId *int   `json:"parentCategoryId"`
}
