package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/clients"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/domain"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/dto"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/google/uuid"
)

//go:generate mockgen -source=./payments_service.go -destination=../test/mock/mock_payments_service.go -package mock
type PaymentService interface {
	GetPayment(id string) (*dto.PaymentResponse, error)
	ProcessPayment(ctx context.Context, req *dto.PaymentRequest) (*dto.PaymentResponse, error)
}

type paymentService struct {
	paymentsRepository repository.PaymentRepository
	bankClient         clients.BankClient
}

var ErrPaymentNotFound = errors.New("payment not found")

func NewPaymentService(repo repository.PaymentRepository, bankClient clients.BankClient) PaymentService {
	return &paymentService{
		paymentsRepository: repo,
		bankClient:         bankClient,
	}
}

func (ps *paymentService) GetPayment(id string) (*dto.PaymentResponse, error) {
	payment := ps.paymentsRepository.GetPayment(id)
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	paymentResponse := dto.FromDomainPayment(payment)
	return paymentResponse, nil
}

func (ps *paymentService) ProcessPayment(ctx context.Context, paymentRequest *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	paymentID := uuid.New().String()

	bankRequest := dto.ToBankRequest(paymentRequest)
	bankResponse, err := ps.bankClient.ProcessPayment(ctx, bankRequest)
	if err != nil {
		if errors.Is(err, clients.ErrBankUnavailable) || errors.Is(err, clients.ErrBankTimeout) {
			payment := paymentRequest.ToDomain(domain.PaymentStatusRejected, paymentID)
			ps.paymentsRepository.AddPayment(payment)
			return dto.FromDomainPayment(payment), nil
		}
		return nil, fmt.Errorf("failed to process payment with bank: %w", err)
	}

	var status domain.PaymentStatus
	if bankResponse.Authorized {
		status = domain.PaymentStatusAuthorized
	} else {
		status = domain.PaymentStatusDeclined
	}

	payment := paymentRequest.ToDomain(status, paymentID)
	ps.paymentsRepository.AddPayment(payment)

	return dto.FromDomainPayment(payment), nil
}
