package order

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Checkout(c *gin.Context) {
	claimsValue, _ := c.Get("user")
	claims := claimsValue.(jwt.MapClaims)

	userID := claims["user_id"].(string)

	order, err := h.service.Checkout(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}

func (h *Handler) Pay(c *gin.Context) {
	orderId := c.Param("order_id")

	err := h.service.Pay(orderId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "paid"})
}
