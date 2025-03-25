package user

import (
	"fmt"
	"net/http"

	configs "github.com/nobregas/ecommerce-mobile-back/config"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/middleware/auth"
	types "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
	router.HandleFunc("/me", auth.WithJwtAuth(h.HandleGetCurrentUser, h.store)).Methods("GET")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, utils.FormatValidationError(err.(validator.ValidationErrors)))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	token, err := auth.CreateJWT([]byte(configs.Envs.JWTSecret), u.ID, u.Role)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, utils.FormatValidationError(err.(validator.ValidationErrors)))
		return
	}

	if _, err := h.store.GetUserByEmail(payload.Email); err == nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("email already registered"))
		return
	}

	if _, err := h.store.GetUserByCPF(payload.Cpf); err == nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("CPF already registered"))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	user := types.User{
		FullName:   payload.FullName,
		Email:      payload.Email,
		Cpf:        payload.Cpf,
		Password:   hashedPassword,
		Role:       types.RoleUser,
		ProfileImg: "https://cdn-icons-png.flaticon.com/512/149/149071.png",
	}

	if err := h.store.CreateUser(user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, user.Sanitize())
}

func (h *Handler) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	userDTO := types.UserDTO{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Cpf:        user.Cpf,
		ProfileImg: user.ProfileImg,
		CreatedAt:  user.CreatedAt,
	}

	utils.WriteJson(w, http.StatusOK, userDTO)
}
