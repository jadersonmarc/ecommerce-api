package cart

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) AddItem(c *gin.Context) {
	var input struct {
		ProductID string
		Quantity  int
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claimsValue, _ := c.Get("user")
	claims := claimsValue.(jwt.MapClaims)

	userID := claims["user_id"].(string)

	cart, err := h.service.AddItems(userID, input.ProductID, input.Quantity)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *Handler) GetCart(c *gin.Context) {
	claimsValue, _ := c.Get("user")
	claims := claimsValue.(jwt.MapClaims)

	userID := claims["user_id"].(string)

	cart, _ := h.service.GetCart(userID)

	c.JSON(200, cart)
}

func (h *Handler) RemoveItem(c *gin.Context) {
	productID := c.Param("product_id")

	claimsValue, _ := c.Get("user")
	claims := claimsValue.(jwt.MapClaims)

	userID := claims["user_id"].(string)

	cart, err := h.service.RemoveItem(userID, productID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, cart)
}

func (h *Handler) UpdateItem(c *gin.Context) {
	productID := c.Param("product_id")

	var input struct {
		Quantity int
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claimsValue, _ := c.Get("user")
	claims := claimsValue.(jwt.MapClaims)

	userID := claims["user_id"].(string)

	cart, err := h.service.UpdateItem(userID, productID, input.Quantity)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, cart)
}
