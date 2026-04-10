package payment

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"
	"github.com/stripe/stripe-go/v85/webhook"
)

type OrderService interface {
	Pay(orderId string) error
}

type Handler struct {
	orderService OrderService
}

func NewHandler(os OrderService) *Handler {
	return &Handler{orderService: os}
}

func (h *Handler) Webhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to read body"})
	}
	event, err := webhook.ConstructEvent(
		body,
		c.GetHeader("Stripe-Signature"),
		"whsec_...",
	)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if event.Type == "payment_intent.succeeded" {

		var pi stripe.PaymentIntent
		json.Unmarshal(event.Data.Raw, &pi)

		orderID := pi.Metadata["order_id"]
		err := h.orderService.Pay(orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	c.JSON(200, gin.H{"status": "ok"})

}
