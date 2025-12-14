package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/dto"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/services"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/validators"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentsHandler struct {
	paymentService services.PaymentService
}

func NewPaymentsHandler(paymentService services.PaymentService) *PaymentsHandler {
	return &PaymentsHandler{
		paymentService: paymentService,
	}
}

// GetPaymentHandler godoc
//
//	@Summary		Get payment details
//	@Description	Retrieve payment information by payment ID
//	@Tags			payments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Payment ID (UUID)"
//	@Success		200	{object}	dto.PaymentResponse
//	@Failure		400	{object}	map[string]string	"Invalid payment ID"
//	@Failure		404	{object}	map[string]string	"Payment not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/api/payments/{id} [get]
func (h *PaymentsHandler) GetPaymentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Header().Set("Content-Type", "application/json")

		if err := uuid.Validate(id); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": errors.New("invalid payment ID").Error()})
			return
		}
		payment, err := h.paymentService.GetPayment(id)

		if err != nil {
			if errors.Is(err, services.ErrPaymentNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(payment); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to encode response"})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}

// ProcessPaymentHandler godoc
//
//	@Summary		Process a payment
//	@Description	Process a new payment transaction
//	@Tags			payments
//	@Accept			json
//	@Produce		json
//	@Param			payment	body		dto.PaymentRequest	true	"Payment request"
//	@Success		200	{object}	dto.PaymentResponse
//	@Failure		400	{object}	dto.ErrorResponse	"Validation failed"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/api/payments [post]
func (h *PaymentsHandler) ProcessPaymentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var paymentRequest dto.PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		/*
			I preferred to use a validator for all fields so that it would be easier for the client to understand
			what was failing in the validation.
		*/
		if validationErrors := validators.ValidatePaymentRequest(paymentRequest); len(validationErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error:  "validation failed",
				Errors: validationErrors,
			})
			return
		}

		paymentResponse, err := h.paymentService.ProcessPayment(r.Context(), &paymentRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(paymentResponse); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to encode response"})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}
