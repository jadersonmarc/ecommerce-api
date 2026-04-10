package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jadersonmarc/ecommerce-api/internal/auth"
	"github.com/jadersonmarc/ecommerce-api/internal/cart"
	"github.com/jadersonmarc/ecommerce-api/internal/order"
	"github.com/jadersonmarc/ecommerce-api/internal/payment"
	"github.com/jadersonmarc/ecommerce-api/internal/product"
	"github.com/jadersonmarc/ecommerce-api/internal/user"
)

func setupRouter() *gin.Engine {
	userRepo := user.NewMemoryRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

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

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/products", productHandler.List)
	r.POST("/webhook", paymentHandler.Webhook)

	authGroup := r.Group("/")
	authGroup.Use(auth.GinAuthMiddleware())
	authGroup.POST("/cart/items", cartHandler.AddItem)
	authGroup.DELETE("/cart/items/:product_id", cartHandler.RemoveItem)
	authGroup.PUT("/cart/items/:product_id", cartHandler.UpdateItem)
	authGroup.GET("/cart", cartHandler.GetCart)
	authGroup.GET("/me", userHandler.Me)
	authGroup.POST("/checkout", orderHandler.Checkout)
	authGroup.POST("/orders/:order_id/pay", orderHandler.Pay)

	adminGroup := r.Group("/")
	adminGroup.Use(auth.GinAuthMiddleware(), auth.GinAdminMiddleware)
	adminGroup.POST("/products", productHandler.Create)

	return r
}

func performJSONRequest(t *testing.T, router http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodeJSONBody(t *testing.T, rec *httptest.ResponseRecorder, target any) {
	t.Helper()

	if err := json.Unmarshal(rec.Body.Bytes(), target); err != nil {
		t.Fatalf("json.Unmarshal() error = %v. body = %s", err, rec.Body.String())
	}
}

func TestAPIRegisterLoginAndGetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	registerRec := performJSONRequest(t, router, http.MethodPost, "/register", map[string]any{
		"name":     "Marc",
		"email":    "marc@test.com",
		"password": "123456",
	}, "")

	if registerRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d. body = %s", http.StatusCreated, registerRec.Code, registerRec.Body.String())
	}

	var registered struct {
		ID       string `json:"ID"`
		Email    string `json:"Email"`
		Password string `json:"Password"`
		Role     string `json:"Role"`
	}
	decodeJSONBody(t, registerRec, &registered)

	if registered.ID == "" {
		t.Fatal("expected registered user ID in response")
	}

	if registered.Email != "marc@test.com" {
		t.Fatalf("expected email %q, got %q", "marc@test.com", registered.Email)
	}

	if registered.Password == "123456" {
		t.Fatal("expected password to be hashed in register response")
	}

	loginRec := performJSONRequest(t, router, http.MethodPost, "/login", map[string]any{
		"email":    "marc@test.com",
		"password": "123456",
	}, "")

	if loginRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d. body = %s", http.StatusOK, loginRec.Code, loginRec.Body.String())
	}

	var loginResponse struct {
		Token string `json:"token"`
	}
	decodeJSONBody(t, loginRec, &loginResponse)

	if loginResponse.Token == "" {
		t.Fatal("expected JWT token in login response")
	}

	meRec := performJSONRequest(t, router, http.MethodGet, "/me", nil, loginResponse.Token)
	if meRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d. body = %s", http.StatusOK, meRec.Code, meRec.Body.String())
	}

	var meResponse struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	decodeJSONBody(t, meRec, &meResponse)

	if meResponse.UserID != registered.ID {
		t.Fatalf("expected user_id %q, got %q", registered.ID, meResponse.UserID)
	}

	if meResponse.Role != "user" {
		t.Fatalf("expected role %q, got %q", "user", meResponse.Role)
	}
}

func TestAPIRejectsAccessToProtectedRouteWithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	rec := performJSONRequest(t, router, http.MethodGet, "/me", nil, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d. body = %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}

	var response struct {
		Error string `json:"error"`
	}
	decodeJSONBody(t, rec, &response)

	if response.Error != "missing token" {
		t.Fatalf("expected error %q, got %q", "missing token", response.Error)
	}
}
