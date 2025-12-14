package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/domain"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/dto"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// mockBankClient is a mock implementation of BankClient for testing
type mockBankClient struct{}

func (m *mockBankClient) ProcessPayment(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error) {
	return &dto.BankPaymentResponse{
		Authorized:        true,
		AuthorizationCode: "test-auth-code",
	}, nil
}

func TestGetPaymentHandler(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	card := domain.Card{
		Number:      "1234567890123456",
		ExpiryMonth: 10,
		ExpiryYear:  2035,
	}
	payment := &domain.Payment{
		ID:       validUUID,
		Card:     card,
		Currency: "GBP",
		Amount:   100,
		Status:   domain.PaymentStatusAuthorized,
	}
	ps := repository.NewPaymentsRepository()
	ps.AddPayment(payment)

	mockBank := &mockBankClient{}
	paymentService := services.NewPaymentService(ps, mockBank)
	payments := NewPaymentsHandler(paymentService)

	r := chi.NewRouter()
	r.Get("/api/payments/{id}", payments.GetPaymentHandler())

	httpServer := &http.Server{
		Addr:    ":8091",
		Handler: r,
	}

	go func() error {
		return httpServer.ListenAndServe()
	}()

	t.Run("PaymentFound", func(t *testing.T) {
		// Create a new HTTP request for testing
		req, _ := http.NewRequest("GET", "/api/payments/"+validUUID, nil)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})
	t.Run("PaymentNotFound", func(t *testing.T) {
		// Create a new HTTP request for testing with a non-existing payment ID
		req, _ := http.NewRequest("GET", "/api/payments/NonExistingID", nil)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the HTTP status code in the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
