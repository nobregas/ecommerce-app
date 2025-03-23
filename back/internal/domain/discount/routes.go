package discount

import (
	"fmt"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	types "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.ProductDiscountStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(
	store types.ProductDiscountStore,
	productStore types.ProductStore,
	userStore types.UserStore) *Handler {

	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	router.HandleFunc("/product/{productID}/discounts",
		auth.WithJwtAuth(h.HandleGetProductDiscounts, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/product/{productID}/discounts/active",
		auth.WithJwtAuth(h.HandleGetActiveDiscounts, h.userStore)).Methods(http.MethodGet)

	// admin routes
	router.HandleFunc("/product/{productID}/discounts",
		auth.WithJwtAuth(auth.WithAdminAuth(h.HandleCreateProductDiscount), h.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/discounts/{discountID}",
		auth.WithJwtAuth(auth.WithAdminAuth(h.HandleUpdateProductDiscount), h.userStore)).Methods(http.MethodPatch)

	router.HandleFunc("/discounts/{discountID}",
		auth.WithJwtAuth(auth.WithAdminAuth(h.HandleDeleteProductDiscount), h.userStore)).Methods(http.MethodDelete)

}

func (h *Handler) HandleGetProductDiscounts(w http.ResponseWriter, r *http.Request) {
	productID := utils.GetParamIdfromPath(r, "productID")

	discounts, err := h.store.GetDiscountsByProduct(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, discounts)
}

func (h *Handler) HandleGetActiveDiscounts(w http.ResponseWriter, r *http.Request) {
	productID := utils.GetParamIdfromPath(r, "productID")

	discounts, err := h.store.GetActiveDiscounts(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, discounts)
}

func (h *Handler) HandleCreateProductDiscount(w http.ResponseWriter, r *http.Request) {
	productID := utils.GetParamIdfromPath(r, "productID")

	var payload *types.CreateProductDiscountPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if _, err := h.productStore.GetProductByID(productID); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	payload.ProductID = productID

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	discount, err := h.store.CreateDiscount(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, discount)
}

func (h *Handler) HandleUpdateProductDiscount(w http.ResponseWriter, r *http.Request) {
	discountID := utils.GetParamIdfromPath(r, "discountID")

	var payload *types.UpdateProductDiscountPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if _, err := h.store.GetDiscoutsByID(discountID); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("discount not found"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	discount, err := h.store.UpdateDiscount(discountID, payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, discount)
}

func (h *Handler) HandleDeleteProductDiscount(w http.ResponseWriter, r *http.Request) {
	discountID := utils.GetParamIdfromPath(r, "discountID")

	if _, err := h.store.GetDiscoutsByID(discountID); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("discount not found"))
		return
	}

	if err := h.store.DeleteDiscount(discountID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusNoContent, nil)
}
