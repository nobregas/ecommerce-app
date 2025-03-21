package rating

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/services/auth"
	"github.com/nobregas/ecommerce-mobile-back/types"
	"github.com/nobregas/ecommerce-mobile-back/utils"
)

type Handler struct {
	store        types.ProductRatingStore
	userStore    types.UserStore
	productStore types.ProductStore
}

func NewHandler(
	store types.ProductRatingStore,
	userStore types.UserStore,
	productStore types.ProductStore) *Handler {

	return &Handler{store: store, userStore: userStore, productStore: productStore}
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
	productID, err := utils.GetParamIdfromPath(r, "productID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// verify if product exists
	_, err = h.productStore.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product with id %d not found", productID))
		return
	}

	// get userID from context
	userID := auth.GetUserIDFromContext(r.Context())

	// get payload from body
	var payload types.CreateProductRatingPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// create rating
	rating, err := h.store.CreateRating(&payload, userID, productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusCreated, rating)
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
