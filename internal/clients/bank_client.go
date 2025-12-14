package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/dto"
)

var (
	ErrBankUnavailable     = errors.New("bank service unavailable")
	ErrBankInvalidRequest  = errors.New("invalid request to bank")
	ErrBankTimeout         = errors.New("bank request timeout")
	ErrBankUnexpectedError = errors.New("unexpected error from bank")
)

//go:generate mockgen -source=./bank_client.go -destination=../test/mock/mock_bank_client.go -package mock
type BankClient interface {
	ProcessPayment(ctx context.Context, req *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error)
}

type bankClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewBankClient(baseURL string) BankClient {
	return &bankClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (bc *bankClient) ProcessPayment(ctx context.Context, bankPaymentRequest *dto.BankPaymentRequest) (*dto.BankPaymentResponse, error) {
	url := fmt.Sprintf("%s/payments", bc.baseURL)

	jsonData, err := json.Marshal(bankPaymentRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bank request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create bank request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := bc.httpClient.Do(httpReq)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrBankTimeout
		}
		return nil, fmt.Errorf("failed to send request to bank: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var bankResp dto.BankPaymentResponse
		if err := json.NewDecoder(resp.Body).Decode(&bankResp); err != nil {
			return nil, fmt.Errorf("failed to decode bank response: %w", err)
		}
		return &bankResp, nil
	case http.StatusBadRequest:
		return nil, ErrBankInvalidRequest
	case http.StatusServiceUnavailable:
		return nil, ErrBankUnavailable
	default:
		return nil, fmt.Errorf("%w: status code %d", ErrBankUnexpectedError, resp.StatusCode)
	}
}
