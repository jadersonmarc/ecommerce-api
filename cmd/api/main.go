package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"

	"github.com/jadersonmarc/ecommerce-api/internal/auth"
	"github.com/jadersonmarc/ecommerce-api/internal/cart"
	"github.com/jadersonmarc/ecommerce-api/internal/database"
	"github.com/jadersonmarc/ecommerce-api/internal/order"
	"github.com/jadersonmarc/ecommerce-api/internal/payment"
	"github.com/jadersonmarc/ecommerce-api/internal/product"
	"github.com/jadersonmarc/ecommerce-api/internal/user"
)

func main() {

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	db, err := database.Open(context.Background(), database.LoadConfig())
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to access sql db from gorm: %v", err)
	}
	defer sqlDB.Close()

	userRepo := user.NewPostgresRepository(db)
	service := user.NewService(userRepo)
	handler := user.NewHandler(service)

	productRepo := product.NewPostgresRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	cartRepo := cart.NewPostgresRepository(db)
	cartService := cart.NewService(cartRepo, productService)
	cartHandler := cart.NewHandler(cartService)

	orderRepo := order.NewPostgresRepository(db)
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
