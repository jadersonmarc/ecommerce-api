package main

import (
	"github.com/gin-gonic/gin"

	"github.com/jadersonmarc/ecommerce-api/internal/auth"
	"github.com/jadersonmarc/ecommerce-api/internal/cart"
	"github.com/jadersonmarc/ecommerce-api/internal/order"
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

	cartRepo := cart.NewMemoryRepository()
	cartService := cart.NewService(cartRepo, productService)
	cartHandler := cart.NewHandler(cartService)

	orderRepo := order.NewMemoryRepository()
	orderService := order.NewService(orderRepo, cartService, productService)
	orderHandler := order.NewHandler(orderService)

	r := gin.Default()

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	r.GET("/products", productHandler.List)

	authGroup := r.Group("/")
	authGroup.Use(auth.GinAuthMiddleware())

	authGroup.POST("/cart/items", cartHandler.AddItem)
	authGroup.POST("/cart/items", cartHandler.AddItem)
	authGroup.DELETE("/cart/items/:product_id", cartHandler.RemoveItem)
	authGroup.PUT("/cart/items/:product_id", cartHandler.UpdateItem)
	authGroup.GET("/cart", cartHandler.GetCart)

	authGroup.GET("/me", handler.Me)

	authGroup.POST("/checkout", orderHandler.Checkout)
	authGroup.POST("/orders/:order_id/pay", orderHandler.Pay)

	adminGroup := r.Group("/")
	adminGroup.Use(auth.GinAdminMiddleware)

	adminGroup.POST("/products", productHandler.Create)

	r.Run(":8080")
}
