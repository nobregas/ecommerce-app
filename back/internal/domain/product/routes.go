package product

import (
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	types "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store          types.ProductStore
	userStore      types.UserStore
	productService types.ProductService
}

func NewHandler(
	store types.ProductStore,
	userStore types.UserStore,
	productService types.ProductService,
) *Handler {
	return &Handler{
		store:          store,
		userStore:      userStore,
		productService: productService,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(auth.WithJwtAuthMiddleware(h.userStore))

	authRouter.HandleFunc("/product/{productID}",
		utils.Compose(
			h.handleGetProductById,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)

	authRouter.HandleFunc("/product",
		utils.Compose(
			h.handleGetProducts,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)

	authRouter.HandleFunc("/product/category/{categoryID}",
		utils.Compose(
			h.handleGetProductsByCategoryID,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)

	authRouter.HandleFunc("/product/details/{productID}",
		utils.Compose(
			h.handleGetProductDetails,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)

	authRouter.HandleFunc("/product/all/details",
		utils.Compose(
			h.handleGetAllSimpleProducts,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)

	// admin routes
	adminRouter := router.PathPrefix("").Subrouter()
	adminRouter.Use(
		auth.WithJwtAuthMiddleware(h.userStore),
		auth.WithAdminAuthMiddleware())

	adminRouter.HandleFunc("/product-with-images",
		utils.Compose(
			h.handleCreateProductWithImages,
			middleware.ErrorHandler,
		)).Methods(http.MethodPost)
}

func (h *Handler) handleGetProductDetails(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	productID := utils.GetParamIdfromPath(r, "productID")

	details := h.productService.GetProductDetails(userID, productID)

	utils.WriteJson(w, http.StatusOK, details)
}

func (h *Handler) handleGetAllSimpleProducts(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	products := h.productService.GetSimpleProducts(userID)

	utils.WriteJson(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProductWithImages(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductWithImagesPayload

	// get json
	if err := utils.ParseJson(r, &payload); err != nil {
		panic(apperrors.NewValidationError("invalid payload", err.Error()))
		return
	}

	createdProduct := h.productService.CreateProductWithImages(payload)

	// response
	utils.WriteJson(w, http.StatusCreated, createdProduct)
}

func (h *Handler) handleUpdateProductByID(w http.ResponseWriter, r *http.Request) {
	// get product id
	productID := utils.GetParamIdfromPath(r, "productID")

	// get payload
	var payload types.UpdateProductPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		panic(apperrors.NewValidationError("invalid payload", err.Error()))
		return
	}

	// update product
	updatedProduct := h.productService.UpdateProductById(productID, payload)

	utils.WriteJson(w, http.StatusOK, updatedProduct)
}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// get product id
	productID := utils.GetParamIdfromPath(r, "productID")

	h.productService.DeleteProduct(productID)

	utils.WriteJson(w, http.StatusNoContent, nil)
}

func (h *Handler) handleGetProductById(w http.ResponseWriter, r *http.Request) {
	productID := utils.GetParamIdfromPath(r, "productID")

	product := h.productService.GetProductByID(productID)
	utils.WriteJson(w, http.StatusOK, product)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products := h.productService.GetProducts()

	utils.WriteJson(w, http.StatusOK, products)
}

func (h *Handler) handleGetProductsByCategoryID(w http.ResponseWriter, r *http.Request) {
	categoryID := utils.GetParamIdfromPath(r, "categoryID")

	products := h.productService.GetProductsByCategoryID(categoryID)

	utils.WriteJson(w, http.StatusOK, products)
}

func (h *Handler) getProductDetails(w http.ResponseWriter, r *http.Request) {
	// get product id from params
	// get product from store
	// response
}

func (h *Handler) getSimpleProducts(w http.ResponseWriter, r *http.Request) {
	// get products from store
	// response
}
