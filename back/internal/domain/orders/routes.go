package orders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
)

type Handler struct {
	orderService types.OrderService
}

func NewHandler(orderService types.OrderService) *Handler {
	return &Handler{
		orderService: orderService,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router, userStore types.UserStore) {
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(auth.WithJwtAuthMiddleware(userStore))

	authRouter.HandleFunc("/orders", h.createOrder).Methods("POST")
	authRouter.HandleFunc("/orders", h.getOrders).Methods("GET")
	authRouter.HandleFunc("/orders/{orderId}", h.getOrderByID).Methods("GET")
	authRouter.HandleFunc("/orders/{orderId}/status", h.updateOrderStatus).Methods("PATCH")
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var payload types.CreateOrderPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := payload.PaymentMethod.Valid(); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid payment method: %s", err.Error())})
		return
	}

	order, err := h.orderService.CreateOrderFromCart(userID, payload.PaymentMethod, payload.PaymentID)
	if err != nil {
		fmt.Printf("[ORDER HANDLER] Error creating order: %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create order: %v", err)})
		return
	}

	utils.WriteJson(w, http.StatusCreated, order)
}

func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	withItems := r.URL.Query().Get("withItems") == "true"

	if withItems {
		orders, err := h.orderService.GetOrdersWithItems(userID)
		if err != nil {
			fmt.Printf("[ORDER HANDLER] Error getting orders with items: %v\n", err)
			utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get orders"})
			return
		}

		if orders == nil {
			orders = []*types.OrderWithItems{}
		}

		utils.WriteJson(w, http.StatusOK, orders)
		return
	}

	orders, err := h.orderService.GetOrdersByUserID(userID)
	if err != nil {
		fmt.Printf("[ORDER HANDLER] Error getting orders: %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get orders"})
		return
	}

	if orders == nil {
		orders = []*types.OrderHistory{}
	}

	utils.WriteJson(w, http.StatusOK, orders)
}

func (h *Handler) getOrderByID(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	orderID := utils.GetParamIdfromPath(r, "orderId")
	if orderID == 0 {
		utils.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
		return
	}

	withItems := r.URL.Query().Get("withItems") == "true"

	if withItems {
		order, err := h.orderService.GetOrderWithItems(orderID)
		if err != nil {
			fmt.Printf("[ORDER HANDLER] Error getting order with items: %v\n", err)
			utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get order"})
			return
		}

		if order.Order.UserID != userID {
			utils.WriteJson(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
			return
		}

		utils.WriteJson(w, http.StatusOK, order)
		return
	}

	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		fmt.Printf("[ORDER HANDLER] Error getting order: %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get order"})
		return
	}

	if order.UserID != userID {
		utils.WriteJson(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
		return
	}

	utils.WriteJson(w, http.StatusOK, order)
}

func (h *Handler) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == 0 {
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	orderID := utils.GetParamIdfromPath(r, "orderId")
	if orderID == 0 {
		utils.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		fmt.Printf("[ORDER HANDLER] Error getting order for status update: %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get order"})
		return
	}

	if order.UserID != userID {
		utils.WriteJson(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
		return
	}

	var payload types.UpdateOrderStatusPayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := payload.Status.Valid(); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid status: %s", err.Error())})
		return
	}

	err = h.orderService.UpdateOrderStatus(orderID, payload.Status)
	if err != nil {
		fmt.Printf("[ORDER HANDLER] Error updating order status: %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update order status"})
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Order status updated successfully"})
}
