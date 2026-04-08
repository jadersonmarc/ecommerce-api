package main

import (
	"github.com/gin-gonic/gin"

	"github.com/jadersonmarc/ecommerce-api/internal/auth"
	"github.com/jadersonmarc/ecommerce-api/internal/product"
	"github.com/jadersonmarc/ecommerce-api/internal/user"
)

func main() {
	repo := user.NewMemoryRepository()
	service := user.NewService(repo)
	handler := user.NewHandler(service)

	productRepo := product.NewMemoryRepository()
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	r := gin.Default()

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	r.GET("/products", productHandler.List)

	authGroup := r.Group("/")
	authGroup.Use(auth.GinAuthMiddleware())

	authGroup.GET("/me", handler.Me)

	adminGroup := r.Group("/")
	adminGroup.Use(auth.GinAdminMiddleware)

	adminGroup.POST("/products", productHandler.Create)

	r.Run(":8080")
}
