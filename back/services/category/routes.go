package category

import (
	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/types"
)

type Handler struct {
	store types.CategoryStore
}

func NewHandler(store types.CategoryStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

}
