package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"

	"github.com/jadersonmarc/ecommerce-api/internal/auth"
	"github.com/jadersonmarc/ecommerce-api/internal/cart"
	"github.com/jadersonmarc/ecommerce-api/internal/order"
	"github.com/jadersonmarc/ecommerce-api/internal/payment"
	"github.com/jadersonmarc/ecommerce-api/internal/product"
	"github.com/jadersonmarc/ecommerce-api/internal/user"
)

func main() {

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

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
	paymentService := payment.NewService()
	orderService := order.NewService(orderRepo, cartService, productService, paymentService)
	orderHandler := order.NewHandler(orderService)
	paymentHandler := payment.NewHandler(orderService)

	r := gin.Default()

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	r.GET("/products", productHandler.List)

	r.POST("/webhook", paymentHandler.Webhook)

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
