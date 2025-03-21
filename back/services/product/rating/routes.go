package rating

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/services/auth"
	"github.com/nobregas/ecommerce-mobile-back/types"
)

type Handler struct {
	store     types.ProductRatingStore
	userStore types.UserStore
}

func NewHandler(store types.ProductRatingStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	router.HandleFunc("/product/{productID}/rating", auth.WithJwtAuth(
		h.HandleCreateProductRating, h.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/product/{productID}/rating", auth.WithJwtAuth(
		h.HandleGetProductRatings, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/user/my/{userID}/rating", auth.WithJwtAuth(
		h.HandleGetMyRatings, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/product/{productID}/rating/average", auth.WithJwtAuth(
		h.HandleGetProductAverageRating, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/product/{productID}/rating", auth.WithJwtAuth(
		h.HandleUpdateProductRating, h.userStore)).Methods(http.MethodPatch)

	router.HandleFunc("/user/my/{userID}/rating", auth.WithJwtAuth(
		h.HandleDeleteMyProductRating, h.userStore)).Methods(http.MethodDelete)

	// admin routes
	router.HandleFunc("/user/{userID}/ratings", auth.WithJwtAuth(
		auth.WithAdminAuth(h.HandleGetUserRatings), h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/product/{productID}/rating", auth.WithJwtAuth(
		auth.WithAdminAuth(h.HandleDeleteProductRating), h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) HandleCreateProductRating(w http.ResponseWriter, r *http.Request) {
	// get product ID from params

	// verify if product exists

	// get userID from context

	// get payload from body

	// validate payload

	// create rating

	// response
}

func (h *Handler) HandleGetProductRatings(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleGetUserRatings(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleGetMyRatings(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleGetProductAverageRating(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleUpdateProductRating(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleDeleteProductRating(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleDeleteMyProductRating(w http.ResponseWriter, r *http.Request) {

}
