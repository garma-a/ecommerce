package main

import (
	"ecom/internal/products"
	"ecom/internal/response"
	"ecom/internal/store"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config
	db *pgxpool.Pool
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

// mount creates the Chi router with middleware and routes configured.
// This follows production-grade patterns with structured logging, recovery, and timeouts.
func (app *application) mount() http.Handler {
	// Initialize structured logger for HTTP requests
	logger := httplog.NewLogger("ecom-api", httplog.Options{
		JSON:             false,
		LogLevel:         slog.LevelInfo,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		Tags: map[string]string{
			"version": "v1.0.0",
			"env":     "production",
		},
	})

	router := chi.NewRouter()

	// Production-grade middleware stack
	router.Use(httplog.RequestLogger(logger))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Heartbeat endpoint for health checks
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{
			"message": "root is fine",
		})
	})

	// Initialize service layer
	queries := store.New(app.db)
	var productService products.IService = products.NewService(queries)
	var productsHandler *products.ProductsHandler = products.NewHandler(productService)

	// Product routes
	router.Get("/products", productsHandler.GetProducts)
	router.Post("/products", productsHandler.CreateProduct)
	router.Get("/products/{id}", productsHandler.GetProductByID)
	router.Delete("/products/{id}", productsHandler.DeleteProduct)

	return router
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	fmt.Printf("Starting server on %s\n", app.config.addr)
	return srv.ListenAndServe()

}
