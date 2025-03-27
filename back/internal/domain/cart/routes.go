package cart

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
)

type Handler struct {
	cartService types.CartService
}

func NewHandler(cartService types.CartService) *Handler {
	return &Handler{
		cartService: cartService,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router, userStore types.UserStore) {

	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(auth.WithJwtAuthMiddleware(userStore))

	authRouter.HandleFunc("/cart", h.createCart).Methods("POST")
	authRouter.HandleFunc("/cart/items", h.getMyCartItems).Methods("GET")
	authRouter.HandleFunc("/cart/items/{productId}", h.addItemToCart).Methods("POST")
	authRouter.HandleFunc("/cart/items/{productId}", h.removeItemFromCart).Methods("DELETE")
	authRouter.HandleFunc("/cart/items/{productId}/remove", h.removeEntireItemFromCart).Methods("DELETE")
	authRouter.HandleFunc("/cart/total", h.getTotal).Methods("GET")
}

func (h *Handler) createCart(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	err := h.cartService.CreateCart(userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create cart"})
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "Cart created successfully"})
}

func (h *Handler) getMyCartItems(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	items, err := h.cartService.GetMyCartItems(userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get cart items"})
		return
	}

	utils.WriteJson(w, http.StatusOK, items)
}

func (h *Handler) addItemToCart(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	productID := utils.GetParamIdfromPath(r, "productId")

	item, err := h.cartService.AddItemToCart(productID, userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add item to cart"})
		return
	}

	utils.WriteJson(w, http.StatusCreated, item)
}

func (h *Handler) removeItemFromCart(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	productID := utils.GetParamIdfromPath(r, "productId")

	err := h.cartService.RemoveItemFromCart(productID, userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove item from cart"})
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Item removed from cart successfully"})
}

func (h *Handler) removeEntireItemFromCart(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	productID := utils.GetParamIdfromPath(r, "productId")

	err := h.cartService.RemoveEntireItemFromCart(productID, userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove item from cart"})
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Item removed from cart successfully"})
}

func (h *Handler) getTotal(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	total, err := h.cartService.GetTotal(userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get cart total"})
		return
	}

	utils.WriteJson(w, http.StatusOK, total)
}
