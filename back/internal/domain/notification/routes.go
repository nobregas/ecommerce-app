package notification

import (
	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
	"net/http"
)

type Handler struct {
	notificationService types.NotificationService
	userStore           types.UserStore
}

func NewHandler(notificationService types.NotificationService, userStore types.UserStore) *Handler {
	return &Handler{
		notificationService: notificationService,
		userStore:           userStore,
	}
}

func (h *Handler) RegisterRouter(router *mux.Router) {
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(auth.WithJwtAuthMiddleware(h.userStore))

	authRouter.HandleFunc("/notification/my",
		utils.Compose(
			h.handleGetNotifications,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)

	authRouter.HandleFunc("/notification/{notificationID}",
		utils.Compose(
			h.handleGetNotificationByID,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)

	adminRouter := router.PathPrefix("").Subrouter()
	adminRouter.Use(
		auth.WithJwtAuthMiddleware(h.userStore),
		auth.WithAdminAuthMiddleware())

	adminRouter.HandleFunc("/notification/{notificationID}",
		utils.Compose(
			h.handleDeleteNotification,
			middleware.ErrorHandler,
		)).Methods(http.MethodDelete)

	adminRouter.HandleFunc("/notification/to/{userID}",
		utils.Compose(
			h.handleCreateNotification,
			middleware.ErrorHandler,
		)).Methods(http.MethodPost)

	adminRouter.HandleFunc("/notification",
		utils.Compose(
			h.handleGetNotifications,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)
}

func (h *Handler) handleGetMyNotifications(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	notifications := h.notificationService.GetMyNotifications(userID)
	utils.WriteJson(w, http.StatusOK, notifications)
}

func (h *Handler) handleGetNotifications(w http.ResponseWriter, r *http.Request) {
	notifications := h.notificationService.GetNotifications()
	utils.WriteJson(w, http.StatusOK, notifications)
}

func (h *Handler) handleGetNotificationByID(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	notificationID := utils.GetParamIdfromPath(r, "notificationID")

	notification := h.notificationService.GetNotificationByID(notificationID)

	if notification.UserID != userID {
		auth.Forbidden(w)
		return
	}
	utils.WriteJson(w, http.StatusOK, notification)
}

func (h *Handler) handleCreateNotification(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetParamIdfromPath(r, "userID")

	var payload types.CreateNotificationPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		panic(apperrors.NewValidationError("invalid payload", err.Error()))
		return
	}

	createdNotification := h.notificationService.CreateNotification(&payload, userID)

	utils.WriteJson(w, http.StatusCreated, createdNotification)
}

func (h *Handler) handleDeleteNotification(w http.ResponseWriter, r *http.Request) {
	notificationID := utils.GetParamIdfromPath(r, "notificationID")

	h.notificationService.DeleteNotification(notificationID)

	utils.WriteJson(w, http.StatusNoContent, nil)
}
