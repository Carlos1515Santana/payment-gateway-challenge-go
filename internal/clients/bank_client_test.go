package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestBankClient_ProcessPayment(t *testing.T) {
	tests := []struct {
		name               string
		mockServerResponse func(w http.ResponseWriter, r *http.Request)
		request            *dto.BankPaymentRequest
		expectedError      error
		expectedAuthorized bool
		expectedAuthCode   string
	}{
		{
			name: "successful payment authorization",
			mockServerResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/payments", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(dto.BankPaymentResponse{
					Authorized:        true,
					AuthorizationCode: "auth-12345",
				})
			},
			request: &dto.BankPaymentRequest{
				CardNumber: "1234567890123456",
				ExpiryDate: "12/2025",
				Currency:   "USD",
				Amount:     1000,
				Cvv:        "123",
			},
			expectedError:      nil,
			expectedAuthorized: true,
			expectedAuthCode:   "auth-12345",
		},
		{
			name: "payment declined by bank",
			mockServerResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(dto.BankPaymentResponse{
					Authorized:        false,
					AuthorizationCode: "",
				})
			},
			request: &dto.BankPaymentRequest{
				CardNumber: "1234567890123456",
				ExpiryDate: "12/2025",
				Currency:   "USD",
				Amount:     1000,
				Cvv:        "123",
			},
			expectedError:      nil,
			expectedAuthorized: false,
			expectedAuthCode:   "",
		},
		{
			name: "bank returns bad request",
			mockServerResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			request: &dto.BankPaymentRequest{
				CardNumber: "invalid",
				ExpiryDate: "12/2025",
				Currency:   "USD",
				Amount:     1000,
				Cvv:        "123",
			},
			expectedError:      ErrBankInvalidRequest,
			expectedAuthorized: false,
			expectedAuthCode:   "",
		},
		{
			name: "bank service unavailable - 503",
			mockServerResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
			},
			request: &dto.BankPaymentRequest{
				CardNumber: "1234567890123456",
				ExpiryDate: "12/2025",
				Currency:   "USD",
				Amount:     1000,
				Cvv:        "123",
			},
			expectedError:      ErrBankUnavailable,
			expectedAuthorized: false,
			expectedAuthCode:   "",
		},
		{
			name: "bank internal server error - 500",
			mockServerResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			request: &dto.BankPaymentRequest{
				CardNumber: "1234567890123456",
				ExpiryDate: "12/2025",
				Currency:   "USD",
				Amount:     1000,
				Cvv:        "123",
			},
			expectedError:      ErrBankUnavailable,
			expectedAuthorized: false,
			expectedAuthCode:   "",
		},
		{
			name: "unexpected status code",
			mockServerResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			},
			request: &dto.BankPaymentRequest{
				CardNumber: "1234567890123456",
				ExpiryDate: "12/2025",
				Currency:   "USD",
				Amount:     1000,
				Cvv:        "123",
			},
			expectedError:      ErrBankUnexpectedError,
			expectedAuthorized: false,
			expectedAuthCode:   "",
		},
		{
			name: "invalid json response from bank",
			mockServerResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			},
			request: &dto.BankPaymentRequest{
				CardNumber: "1234567890123456",
				ExpiryDate: "12/2025",
				Currency:   "USD",
				Amount:     1000,
				Cvv:        "123",
			},
			expectedError:      nil,
			expectedAuthorized: false,
			expectedAuthCode:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.mockServerResponse))
			defer server.Close()

			client := NewBankClient(server.URL)
			result, err := client.ProcessPayment(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else if tt.name == "invalid json response from bank" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to decode bank response")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedAuthorized, result.Authorized)
				assert.Equal(t, tt.expectedAuthCode, result.AuthorizationCode)
			}
		})
	}
}
