package payment

import (
	"github.com/stripe/stripe-go/v85"
	"github.com/stripe/stripe-go/v85/paymentintent"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreatePaymentIntent(amount int64, orderID string) (string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String("blr"),
	}

	params.Metadata = map[string]string{
		"order_id": orderID,
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", err
	}

	return pi.ClientSecret, nil
}
