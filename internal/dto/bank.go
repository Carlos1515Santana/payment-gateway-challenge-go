package dto

import "fmt"

type BankPaymentRequest struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	Currency   string `json:"currency"`
	Amount     int    `json:"amount"`
	Cvv        string `json:"cvv"`
}

type BankPaymentResponse struct {
	Authorized        bool   `json:"authorized"`
	AuthorizationCode string `json:"authorization_code"`
}

func ToBankRequest(req *PaymentRequest) *BankPaymentRequest {
	expiryDate := formatExpiryDate(req.ExpiryMonth, req.ExpiryYear)

	return &BankPaymentRequest{
		CardNumber: req.CardNumber,
		ExpiryDate: expiryDate,
		Currency:   req.Currency,
		Amount:     req.Amount,
		Cvv:        req.Cvv,
	}
}

func formatExpiryDate(month, year int) string {
	return fmt.Sprintf("%02d/%d", month, year)
}
