package category

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
	store     types.CategoryStore
	userStore types.UserStore
}

func NewHandler(store types.CategoryStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	router.HandleFunc("/category", auth.WithJwtAuth(h.handleGetCategories, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/category/{categoryID}", auth.WithJwtAuth(
		h.handleGetCategoryByID, h.userStore)).Methods(http.MethodGet)

	// admin routes
	router.HandleFunc("/category", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleCreateCategory), h.userStore)).Methods(http.MethodPost)

	router.HandleFunc("/category/{categoryID}", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleUpdateCategoryByID), h.userStore)).Methods(http.MethodPatch)

	router.HandleFunc("/category/{categoryID}", auth.WithJwtAuth(
		auth.WithAdminAuth(h.handleDeleteCategory), h.userStore)).Methods(http.MethodDelete)

}

func (h *Handler) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.GetCategories()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, categories)
}

func (h *Handler) handleGetCategoryByID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := utils.GetParamIdfromPath(r, "categoryID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	category, err := h.store.GetCategoryByID(categoryID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("category with id %d not found", categoryID))
		return
	}

	utils.WriteJson(w, http.StatusOK, category)
}

func (h *Handler) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateCategoryPayload

	// get json
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// verify if the parent category exists
	if payload.ParentCategoryId != nil {
		_, err := h.store.GetCategoryByID(*payload.ParentCategoryId)
		if err != nil {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("parent category with id %d not found", *payload.ParentCategoryId))
			return
		}
	}

	// create category
	createdCategory, err := h.store.CreateCategory(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, createdCategory)
}

func (h *Handler) handleUpdateCategoryByID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := utils.GetParamIdfromPath(r, "categoryID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// verify if the category exists
	_, err = h.store.GetCategoryByID(categoryID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("category with id %d not found", categoryID))
		return
	}

	// get payload
	var payload types.UpdateCategoryPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// update
	updatedCategory, err := h.store.UpdateCategory(categoryID, payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, updatedCategory)
}

func (h *Handler) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := utils.GetParamIdfromPath(r, "categoryID")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// verify if the category exists
	_, err = h.store.GetCategoryByID(categoryID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("category with id %d not found", categoryID))
		return
	}

	// delete category
	if err := h.store.DeleteCategory(categoryID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusNoContent, nil)
}
