package utils

import "net/http"

func Compose(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ex of use:
//router.HandleFunc("/product",
//	Compose(
//		h.handleCreateProduct,
//		auth.WithAdminAuth,
//		auth.WithJwtAuth,
//		ErrorHandler,
//	),
//).Methods(http.MethodPost)
