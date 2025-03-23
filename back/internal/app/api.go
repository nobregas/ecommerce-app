package app

import (
	"database/sql"
	category "github.com/nobregas/ecommerce-mobile-back/internal/domain/category"
	"github.com/nobregas/ecommerce-mobile-back/internal/domain/discount"
	product "github.com/nobregas/ecommerce-mobile-back/internal/domain/product"
	"github.com/nobregas/ecommerce-mobile-back/internal/domain/rating"
	user "github.com/nobregas/ecommerce-mobile-back/internal/domain/user"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

	userStore := user.NewStore(s.db)
	productStore := product.NewStore(s.db)
	categoryStore := category.NewStore(s.db)
	discountStore := discount.NewStore(s.db)
	ratingStore := rating.NewStore(s.db)

	productService := product.NewProductService(
		productStore,
		userStore,
		discountStore,
		ratingStore)

	// user
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// product
	productHandler := product.NewHandler(productStore, userStore, productService)
	productHandler.RegisterRoutes(subrouter)

	// category
	categoryHandler := category.NewHandler(categoryStore, userStore)
	categoryHandler.RegisterRoutes(subrouter)

	// discount
	discountHandler := discount.NewHandler(discountStore, productStore, userStore)
	discountHandler.RegisterRoutes(subrouter)

	// rating
	ratingHandler := rating.NewHandler(ratingStore, userStore, productStore)
	ratingHandler.RegisterRoutes(subrouter)

	log.Printf("Server listening on %s", s.addr)

	return http.ListenAndServe(s.addr, router)
}
