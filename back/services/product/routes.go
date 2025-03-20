package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/services/auth"
	"github.com/nobregas/ecommerce-mobile-back/types"
	"github.com/nobregas/ecommerce-mobile-back/utils"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	router.HandleFunc("/product", auth.WithJwtAuth(h.handleGetProducts, h.userStore)).Methods(http.MethodGet)

	// admin routes
	router.HandleFunc("/product", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleCreateProduct), h.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/product-with-images", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleCreateProductWithImages), h.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/product/{productID}", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleUpdateProductByID), h.userStore)).Methods(http.MethodPatch)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, ps)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductPayload

	// get json
	if err := utils.ParseJson(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate
	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// create product
	err := h.store.CreateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusCreated, product)
}

func (h *Handler) handleCreateProductWithImages(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductWithImagesPayload

	// get json
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// create product
	createdProduct, err := h.store.CreateProductWithImages(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	utils.WriteJson(w, http.StatusCreated, createdProduct)
}

func (h *Handler) handleUpdateProductByID(w http.ResponseWriter, r *http.Request) {
	// get product id
	vars := mux.Vars(r)
	str, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing productID"))
		return
	}

	// convert product id to str
	productID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid productID"))
		return
	}

	// verify if the product exists
	_, err = h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product with id %d not found", productID))
		return
	}

	// get payload
	var payload types.UpdateProductPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// update
	if err := h.store.UpdateProduct(productID, payload); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// response
	updatedProduct, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, updatedProduct)
}
