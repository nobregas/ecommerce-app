package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nobregas/ecommerce-mobile-back/services/category"
	"github.com/nobregas/ecommerce-mobile-back/services/product"
	"github.com/nobregas/ecommerce-mobile-back/services/product/discount"
	"github.com/nobregas/ecommerce-mobile-back/services/product/rating"
	"github.com/nobregas/ecommerce-mobile-back/services/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// user
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// product
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	// category
	categoryStore := category.NewStore(s.db)
	categoryHandler := category.NewHandler(categoryStore, userStore)
	categoryHandler.RegisterRoutes(subrouter)

	// discount
	discountStore := discount.NewStore(s.db)
	discountHandler := discount.NewHandler(discountStore, productStore, userStore)
	discountHandler.RegisterRoutes(subrouter)

	// rating
	ratingStore := rating.NewStore(s.db)
	ratingHandler := rating.NewHandler(ratingStore, userStore, productStore)
	ratingHandler.RegisterRoutes(subrouter)

	log.Printf("Server listening on %s", s.addr)

	return http.ListenAndServe(s.addr, router)
}
