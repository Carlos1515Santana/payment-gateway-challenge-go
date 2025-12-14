package api

import (
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/clients"
	handler "github.com/cko-recruitment/payment-gateway-challenge-go/internal/handlers"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/services"
)

type handlers struct {
	PaymentsHandler *handler.PaymentsHandler
	SwaggerHandler  *handler.SwaggerHandler
	PingHandler     *handler.PingHandler
}

func loadDependencies() *handlers {
	paymentsRepo := repository.NewPaymentsRepository()
	// Could be use a config file to store the bank client URL
	bankClient := clients.NewBankClient("http://localhost:8080")
	paymentsService := services.NewPaymentService(paymentsRepo, bankClient)

	return &handlers{
		PaymentsHandler: handler.NewPaymentsHandler(paymentsService),
		SwaggerHandler:  handler.NewSwaggerHandler(),
		PingHandler:     handler.NewPingHandler(),
	}
}
