package dto

import "github.com/cko-recruitment/payment-gateway-challenge-go/internal/domain"

// The cvv and card number could be encrypted or not stored, depends on the security policy.
type PaymentRequest struct {
	CardNumber  string `json:"card_number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	Currency    string `json:"currency"`
	Amount      int    `json:"amount"`
	Cvv         string `json:"cvv"`
}

type PaymentResponse struct {
	ID                 string `json:"id"`
	PaymentStatus      string `json:"payment_status"`
	CardNumberLastFour string `json:"card_number_last_four"`
	ExpiryMonth        int    `json:"expiry_month"`
	ExpiryYear         int    `json:"expiry_year"`
	Currency           string `json:"currency"`
	Amount             int    `json:"amount"`
}

func (pr *PaymentRequest) ToDomain(status domain.PaymentStatus, paymentID string) *domain.Payment {
	card := domain.Card{
		Number:      pr.CardNumber,
		ExpiryMonth: pr.ExpiryMonth,
		ExpiryYear:  pr.ExpiryYear,
	}

	return domain.NewPayment(paymentID, card, pr.Currency, pr.Amount, status)
}

func FromDomainPayment(payment *domain.Payment) *PaymentResponse {
	lastFour := payment.Card.GetLastFourDigits()

	return &PaymentResponse{
		ID:                 payment.ID,
		PaymentStatus:      string(payment.Status),
		CardNumberLastFour: lastFour,
		ExpiryMonth:        payment.Card.ExpiryMonth,
		ExpiryYear:         payment.Card.ExpiryYear,
		Currency:           payment.Currency,
		Amount:             payment.Amount,
	}
}
