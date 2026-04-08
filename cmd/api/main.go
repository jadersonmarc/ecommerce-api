package main

import (
	"log"
	"net/http"

	"github.com/jadersonmarc/ecommerce-api/internal/auth"
	"github.com/jadersonmarc/ecommerce-api/internal/user"
)

func main() {
	repo := user.NewMemoryRepository()
	service := user.NewService(repo)
	handler := user.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/login", handler.Login)

	mux.Handle("/me", auth.AuthMiddleware(http.HandlerFunc(handler.Me)))

	log.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
