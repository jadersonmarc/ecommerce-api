package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c *gin.Context) {
	var input struct {
		Name        string
		Description string
		Price       int64
		Stock       int
	}

	c.ShouldBindJSON(&input)

	p, err := h.service.Create(
		input.Name,
		input.Description,
		input.Price,
		input.Stock,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *Handler) List(c *gin.Context) {
	products, err := h.service.List()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, products)
}
