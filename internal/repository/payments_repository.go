package repository

import (
	"sync"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/domain"
)

//go:generate mockgen -source=./payments_repository.go -destination=../test/mock/mock_payments_repository.go -package mock
type PaymentRepository interface {
	GetPayment(id string) *domain.Payment
	AddPayment(payment *domain.Payment)
}

type paymentRepository struct {
	mu       sync.RWMutex
	payments map[string]*domain.Payment
}

func NewPaymentsRepository() PaymentRepository {
	return &paymentRepository{
		payments: make(map[string]*domain.Payment),
	}
}

func (ps *paymentRepository) GetPayment(id string) *domain.Payment {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if payment, exists := ps.payments[id]; exists {
		return payment
	}
	return nil
}

func (ps *paymentRepository) AddPayment(payment *domain.Payment) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.payments[payment.ID] = payment
}
