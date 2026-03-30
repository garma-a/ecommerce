package main

import (
	"ecom/internal/products"
	"ecom/internal/store"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
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

// create the router using Gin and return it back
func (app *application) mount() http.Handler {

	router := gin.New()
	router.Use(requestid.New())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(timeout.New(
		timeout.WithTimeout(60 * time.Second),
	))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "root is fine",
		})

	})
	queries := store.New(app.db)
	var productService products.IService = products.NewService(queries)
	var productsHandler *products.ProductsHandler = products.NewHandler(productService)
	router.GET("/products", productsHandler.GetProducts)
	router.GET("/products/:id", productsHandler.GetProductByID)
	router.POST("/products", productsHandler.CreateProduct)
	router.DELETE("/products/:id", productsHandler.DeleteProduct)

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
