package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/nobregas/ecommerce-mobile-back/internal/domain/cart"
	category "github.com/nobregas/ecommerce-mobile-back/internal/domain/category"
	"github.com/nobregas/ecommerce-mobile-back/internal/domain/discount"
	"github.com/nobregas/ecommerce-mobile-back/internal/domain/favorite"
	"github.com/nobregas/ecommerce-mobile-back/internal/domain/notification"
	"github.com/nobregas/ecommerce-mobile-back/internal/domain/orders"
	product "github.com/nobregas/ecommerce-mobile-back/internal/domain/product"
	"github.com/nobregas/ecommerce-mobile-back/internal/domain/rating"
	user "github.com/nobregas/ecommerce-mobile-back/internal/domain/user"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	productStore := product.NewStore(s.db)
	categoryStore := category.NewStore(s.db)
	discountStore := discount.NewStore(s.db)
	ratingStore := rating.NewStore(s.db)
	notificationStore := notification.NewStore(s.db)
	favoriteStore := favorite.NewStore(s.db)
	cartStore := cart.NewStore(s.db)
	orderStore := orders.NewStore(s.db)

	cartService := cart.NewService(
		cartStore,
		productStore,
		discountStore,
	)

	orderService := orders.NewService(
		orderStore,
		cartStore,
		productStore,
	)
	userStore := user.NewStore(s.db, cartService)

	productService := product.NewProductService(
		productStore,
		userStore,
		discountStore,
		ratingStore,
	)

	favoriteService := favorite.NewService(
		favoriteStore,
		userStore,
		productStore)

	notificationService := notification.NewNotificationService(notificationStore, userStore)

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

	// notification
	notificationHandler := notification.NewHandler(notificationService, userStore)
	notificationHandler.RegisterRouter(subrouter)

	// favorite
	favoriteHandler := favorite.NewHandler(favoriteService, userStore)
	favoriteHandler.RegisterRouter(subrouter)

	// cart
	cartHandler := cart.NewHandler(cartService)
	cartHandler.RegisterRoutes(subrouter, userStore)

	// order
	orderHandler := orders.NewHandler(orderService)
	orderHandler.RegisterRoutes(subrouter, userStore)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	log.Printf("Server listening on %s", s.addr)

	return http.ListenAndServe(s.addr, handler)
}
