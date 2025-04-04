package rating

import (
	"fmt"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	types "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
	"net/http"

	"github.com/gorilla/mux"
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

	router.HandleFunc("/user/my/rating", auth.WithJwtAuth(
		h.HandleGetMyRatings, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/product/{productID}/rating/average", auth.WithJwtAuth(
		h.HandleGetProductAverageRating, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/product/{productID}/rating", auth.WithJwtAuth(
		h.HandleUpdateProductRating, h.userStore)).Methods(http.MethodPatch)

	router.HandleFunc("/user/my/rating/{ratingID}", auth.WithJwtAuth(
		h.HandleDeleteMyProductRating, h.userStore)).Methods(http.MethodDelete)

	// admin routes
	router.HandleFunc("/user/{userID}/rating", auth.WithJwtAuth(
		auth.WithAdminAuth(h.HandleGetUserRatings), h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/product/{ratingID}/rating", auth.WithJwtAuth(
		auth.WithAdminAuth(h.HandleDeleteProductRating), h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) HandleCreateProductRating(w http.ResponseWriter, r *http.Request) {
	// get product ID from params
	productID := utils.GetParamIdfromPath(r, "productID")

	// verify if product exists
	_, err := h.productStore.GetProductByID(productID)
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
	// get product ID from params
	productID := utils.GetParamIdfromPath(r, "productID")

	// verify if product exists
	_, err := h.productStore.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product with id %d not found", productID))
		return
	}

	// get ratings
	ratings, err := h.store.GetRatingsByProduct(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusOK, ratings)
}

func (h *Handler) HandleGetUserRatings(w http.ResponseWriter, r *http.Request) {
	// get user ID from params
	userID := utils.GetParamIdfromPath(r, "userID")

	// verify if user exists
	_, err := h.userStore.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user with id %d not found", userID))
		return
	}

	// get ratings
	ratings, err := h.store.GetRatingsByUser(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusOK, ratings)
}

func (h *Handler) HandleGetMyRatings(w http.ResponseWriter, r *http.Request) {
	// get user ID from context
	userID := auth.GetUserIDFromContext(r.Context())

	// get ratings
	ratings, err := h.store.GetRatingsByUser(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusOK, ratings)
}

func (h *Handler) HandleGetProductAverageRating(w http.ResponseWriter, r *http.Request) {
	// get product ID from params
	productID := utils.GetParamIdfromPath(r, "productID")

	// verify if product exists
	_, err := h.productStore.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product with id %d not found", productID))
		return
	}

	// get average rating
	averageRating, err := h.store.GetAverageRating(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusOK, averageRating)
}

func (h *Handler) HandleUpdateProductRating(w http.ResponseWriter, r *http.Request) {
	// get rating ID from params
	ratingID := utils.GetParamIdfromPath(r, "ratingID")

	// verify if rating exists
	rating, err := h.store.GetRating(ratingID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("rating with id %d not found", ratingID))
		return
	}

	// get userID from context
	userID := auth.GetUserIDFromContext(r.Context())

	// verify if rating belongs to user
	if rating.UserID != userID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("rating with id %d does not belong to user with id %d", ratingID, userID))
		return
	}

	// get payload from body
	var payload types.UpdateProductRatingPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// update rating
	updatedRating, err := h.store.UpdateRating(ratingID, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusOK, updatedRating)
}

func (h *Handler) HandleDeleteProductRating(w http.ResponseWriter, r *http.Request) {
	// get rating ID from params
	ratingID := utils.GetParamIdfromPath(r, "ratingID")

	// verify if rating exists
	_, err := h.store.GetRating(ratingID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("rating with id %d not found", ratingID))
		return
	}

	// delete rating
	if err := h.store.DeleteRating(ratingID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusNoContent, nil)
}

func (h *Handler) HandleDeleteMyProductRating(w http.ResponseWriter, r *http.Request) {
	// get rating ID from params
	ratingID := utils.GetParamIdfromPath(r, "ratingID")

	// verify if rating exists
	rating, err := h.store.GetRating(ratingID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("rating with id %d not found", ratingID))
		return
	}

	// get user ID from context
	userID := auth.GetUserIDFromContext(r.Context())

	// verify if rating belongs to user
	if rating.UserID != userID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("rating with id %d does not belong to user with id %d", ratingID, userID))
		return
	}

	// delete rating
	if err := h.store.DeleteRating(ratingID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusNoContent, nil)
}
