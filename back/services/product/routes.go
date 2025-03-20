package product

import (
	"fmt"
	"net/http"

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

	router.HandleFunc("/product/{productID}", auth.WithJwtAuth(
		h.handleGetProductByID, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/product/category/{categoryID}", auth.WithJwtAuth(
		h.handleGetProductsByCategoryID, h.userStore)).Methods(http.MethodGet)

	// admin routes
	router.HandleFunc("/product", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleCreateProduct), h.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/product-with-images", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleCreateProductWithImages), h.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/product/{productID}", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleUpdateProductByID), h.userStore)).Methods(http.MethodPatch)

	router.HandleFunc("/product/{productID}", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleDeleteProduct), h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, ps)
}

func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	productID, err := utils.GetParamIdfromPath(r, "productID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product with id %d not found", productID))
		return
	}

	utils.WriteJson(w, http.StatusOK, product)
}

func (h *Handler) handleGetProductsByCategoryID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := utils.GetParamIdfromPath(r, "categoryID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	products, err := h.store.GetProductsByCategory(categoryID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, products)
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
	productID, err := utils.GetParamIdfromPath(r, "productID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
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

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// get product id
	productID, err := utils.GetParamIdfromPath(r, "productID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// verify if product exists
	_, err = h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product with id %d not found", productID))
		return
	}

	// delete product
	if err := h.store.DeleteProduct(productID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusNoContent, nil)
}
