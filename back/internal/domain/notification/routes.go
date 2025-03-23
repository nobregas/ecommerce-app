package notification

import "github.com/gorilla/mux"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRouter(router *mux.Router) {

}
