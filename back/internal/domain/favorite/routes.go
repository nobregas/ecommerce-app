package favorite

import (
	"github.com/gorilla/mux"
	_ "github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
	"net/http"
)

type Handler struct {
	favoriteService types.UserFavoriteService
	userStore       types.UserStore
}

func NewHandler(favoriteService types.UserFavoriteService, userStore types.UserStore) *Handler {
	return &Handler{
		favoriteService: favoriteService,
		userStore:       userStore,
	}
}

func (h *Handler) RegisterRouter(router *mux.Router) {
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(auth.WithJwtAuthMiddleware(h.userStore))

	authRouter.HandleFunc("/favorite",
		utils.Compose(
			h.handleGetUserFavorites,
			middleware.ErrorHandler,
		)).Methods(http.MethodGet)

	authRouter.HandleFunc("/favorite/{productID}",
		utils.Compose(
			h.handleAddFavorite,
			middleware.ErrorHandler,
		)).Methods(http.MethodPost)

	authRouter.HandleFunc("/favorite/{productID}",
		utils.Compose(
			h.handleRemoveFavorite,
			middleware.ErrorHandler,
		)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetUserFavorites(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	favorites := h.favoriteService.GetUserFavorite(userID)
	utils.WriteJson(w, http.StatusOK, favorites)
}

func (h *Handler) handleAddFavorite(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	productID := utils.GetParamIdfromPath(r, "productID")

	favorite := h.favoriteService.AddFavorite(userID, productID)
	utils.WriteJson(w, http.StatusCreated, favorite)
}

func (h *Handler) handleRemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	productID := utils.GetParamIdfromPath(r, "productID")

	h.favoriteService.RemoveFavorite(userID, productID)
	w.WriteHeader(http.StatusNoContent)
}
