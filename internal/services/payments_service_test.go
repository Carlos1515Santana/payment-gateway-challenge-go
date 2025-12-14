package services

import (
	"context"
	"errors"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/clients"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/domain"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/dto"
	"github.com/stretchr/testify/assert"
)

type mockPaymentRepository struct {
	getPaymentFunc func(id string) *domain.Payment
	addPaymentFunc func(payment *domain.Payment)
}

func (m *mockPaymentRepository) GetPayment(id string) *domain.Payment {
	if m.getPaymentFunc != nil {
		return m.getPaymentFunc(id)
	}
	return nil
}

func (m *mockPaymentRepository) AddPayment(payment *domain.Payment) {
	if m.addPaymentFunc != nil {
		m.addPaymentFunc(payment)
	}
}

type mockBankClient struct {
	processPaymentFunc func(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error)
}

func (m *mockBankClient) ProcessPayment(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error) {
	if m.processPaymentFunc != nil {
		return m.processPaymentFunc(ctx, req)
	}
	return nil, nil
}

func TestPaymentService_GetPayment(t *testing.T) {
	tests := []struct {
		name          string
		paymentID     string
		mockRepo      *mockPaymentRepository
		expectedError error
		expectedID    string
	}{
		{
			name:      "payment found successfully",
			paymentID: "test-id-123",
			mockRepo: &mockPaymentRepository{
				getPaymentFunc: func(id string) *domain.Payment {
					if id == "test-id-123" {
						return &domain.Payment{
							ID: "test-id-123",
							Card: domain.Card{
								Number:      "1234567890123456",
								ExpiryMonth: 12,
								ExpiryYear:  2025,
							},
							Currency: "USD",
							Amount:   1000,
							Status:   domain.PaymentStatusAuthorized,
						}
					}
					return nil
				},
			},
			expectedError: nil,
			expectedID:    "test-id-123",
		},
		{
			name:      "payment not found",
			paymentID: "non-existent-id",
			mockRepo: &mockPaymentRepository{
				getPaymentFunc: func(id string) *domain.Payment {
					return nil
				},
			},
			expectedError: ErrPaymentNotFound,
			expectedID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBank := &mockBankClient{}

			service := NewPaymentService(tt.mockRepo, mockBank)
			result, err := service.GetPayment(tt.paymentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedID, result.ID)
			}
		})
	}
}

func TestPaymentService_ProcessPayment(t *testing.T) {
	tests := []struct {
		name           string
		paymentRequest *dto.PaymentRequest
		mockRepo       *mockPaymentRepository
		mockBank       *mockBankClient
		expectedStatus domain.PaymentStatus
		expectedError  error
	}{
		{
			name: "payment authorized successfully",
			paymentRequest: &dto.PaymentRequest{
				CardNumber:  "1234567890123456",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			},
			mockRepo: &mockPaymentRepository{
				addPaymentFunc: func(payment *domain.Payment) {
					assert.Equal(t, domain.PaymentStatusAuthorized, payment.Status)
				},
			},
			mockBank: &mockBankClient{
				processPaymentFunc: func(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error) {
					return &dto.BankPaymentResponse{
						Authorized:        true,
						AuthorizationCode: "auth-123",
					}, nil
				},
			},
			expectedStatus: domain.PaymentStatusAuthorized,
			expectedError:  nil,
		},
		{
			name: "payment declined by bank",
			paymentRequest: &dto.PaymentRequest{
				CardNumber:  "1234567890123456",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			},
			mockRepo: &mockPaymentRepository{
				addPaymentFunc: func(payment *domain.Payment) {
					assert.Equal(t, domain.PaymentStatusDeclined, payment.Status)
				},
			},
			mockBank: &mockBankClient{
				processPaymentFunc: func(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error) {
					return &dto.BankPaymentResponse{
						Authorized:        false,
						AuthorizationCode: "",
					}, nil
				},
			},
			expectedStatus: domain.PaymentStatusDeclined,
			expectedError:  nil,
		},
		{
			name: "bank unavailable - payment rejected",
			paymentRequest: &dto.PaymentRequest{
				CardNumber:  "1234567890123456",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			},
			mockRepo: &mockPaymentRepository{
				addPaymentFunc: func(payment *domain.Payment) {
					assert.Equal(t, domain.PaymentStatusRejected, payment.Status)
				},
			},
			mockBank: &mockBankClient{
				processPaymentFunc: func(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error) {
					return nil, clients.ErrBankUnavailable
				},
			},
			expectedStatus: domain.PaymentStatusRejected,
			expectedError:  nil,
		},
		{
			name: "bank timeout - payment rejected",
			paymentRequest: &dto.PaymentRequest{
				CardNumber:  "1234567890123456",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			},
			mockRepo: &mockPaymentRepository{
				addPaymentFunc: func(payment *domain.Payment) {
					assert.Equal(t, domain.PaymentStatusRejected, payment.Status)
				},
			},
			mockBank: &mockBankClient{
				processPaymentFunc: func(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error) {
					return nil, clients.ErrBankTimeout
				},
			},
			expectedStatus: domain.PaymentStatusRejected,
			expectedError:  nil,
		},
		{
			name: "unexpected bank error",
			paymentRequest: &dto.PaymentRequest{
				CardNumber:  "1234567890123456",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			},
			mockRepo: &mockPaymentRepository{},
			mockBank: &mockBankClient{
				processPaymentFunc: func(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error) {
					return nil, errors.New("unexpected error")
				},
			},
			expectedStatus: "",
			expectedError:  errors.New("failed to process payment with bank"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewPaymentService(tt.mockRepo, tt.mockBank)
			result, err := service.ProcessPayment(context.Background(), tt.paymentRequest)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, string(tt.expectedStatus), result.PaymentStatus)
				assert.NotEmpty(t, result.ID)
			}
		})
	}
}
