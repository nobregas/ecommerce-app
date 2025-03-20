package types

type CategoryStore interface {
}

type Category struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ImageUrl         string `json:"imageUrl"`
	ParentCategoryId *int   `json:"parentCategoryId"`
}
