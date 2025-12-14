package validators

import (
	"testing"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/dto"
)

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name          string
		amount        int
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "valid positive amount",
			amount:        1000,
			expectedError: false,
		},
		{
			name:          "valid amount of 1",
			amount:        1,
			expectedError: false,
		},
		{
			name:          "invalid zero amount",
			amount:        0,
			expectedError: true,
			errorMessage:  "amount must be a positive integer",
		},
		{
			name:          "invalid negative amount",
			amount:        -100,
			expectedError: true,
			errorMessage:  "amount must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateAmount(tt.amount)

			if tt.expectedError {
				if len(errors) == 0 {
					t.Errorf("expected error but got none")
				} else {
					if errors[0].Field != "amount" {
						t.Errorf("expected field 'amount', got '%s'", errors[0].Field)
					}
					if errors[0].Message != tt.errorMessage {
						t.Errorf("expected message '%s', got '%s'", tt.errorMessage, errors[0].Message)
					}
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("expected no errors but got: %v", errors)
				}
			}
		})
	}
}

func TestValidateCardNumber(t *testing.T) {
	tests := []struct {
		name           string
		cardNumber     string
		expectedErrors int
		errorMessages  []string
	}{
		{
			name:           "valid 16 digit card",
			cardNumber:     "4532015112830366",
			expectedErrors: 0,
		},
		{
			name:           "valid 14 digit card",
			cardNumber:     "12345678901234",
			expectedErrors: 0,
		},
		{
			name:           "valid 19 digit card",
			cardNumber:     "1234567890123456789",
			expectedErrors: 0,
		},
		{
			name:           "empty card number",
			cardNumber:     "",
			expectedErrors: 1,
			errorMessages:  []string{"card number is required"},
		},
		{
			name:           "card number too short",
			cardNumber:     "1234567890123",
			expectedErrors: 1,
			errorMessages:  []string{"card number must be between 14-19 characters"},
		},
		{
			name:           "card number too long",
			cardNumber:     "12345678901234567890",
			expectedErrors: 1,
			errorMessages:  []string{"card number must be between 14-19 characters"},
		},
		{
			name:           "card number with letters",
			cardNumber:     "4532015112830abc",
			expectedErrors: 1,
			errorMessages:  []string{"card number must contain only numeric characters"},
		},
		{
			name:           "card number with spaces",
			cardNumber:     "4532 0151 1283 0366",
			expectedErrors: 1,
			errorMessages:  []string{"card number must contain only numeric characters"},
		},
		{
			name:           "card number with special characters",
			cardNumber:     "4532-0151-1283-0366",
			expectedErrors: 1,
			errorMessages:  []string{"card number must contain only numeric characters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateCardNumber(tt.cardNumber)

			if len(errors) != tt.expectedErrors {
				t.Errorf("expected %d errors, got %d: %v", tt.expectedErrors, len(errors), errors)
			}

			for i, msg := range tt.errorMessages {
				if i < len(errors) && errors[i].Message != msg {
					t.Errorf("expected error message '%s' at index %d, got '%s'", msg, i, errors[i].Message)
				}
			}
		})
	}
}

func TestValidateExpiryMonth(t *testing.T) {
	tests := []struct {
		name          string
		month         int
		expectedError bool
	}{
		{
			name:          "valid month 1",
			month:         1,
			expectedError: false,
		},
		{
			name:          "valid month 12",
			month:         12,
			expectedError: false,
		},
		{
			name:          "valid month 6",
			month:         6,
			expectedError: false,
		},
		{
			name:          "invalid month 0",
			month:         0,
			expectedError: true,
		},
		{
			name:          "invalid month 13",
			month:         13,
			expectedError: true,
		},
		{
			name:          "invalid negative month",
			month:         -1,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateExpiryMonth(tt.month)

			if tt.expectedError && len(errors) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.expectedError && len(errors) > 0 {
				t.Errorf("expected no errors but got: %v", errors)
			}
			if tt.expectedError && len(errors) > 0 && errors[0].Field != "expiry_month" {
				t.Errorf("expected field 'expiry_month', got '%s'", errors[0].Field)
			}
		})
	}
}

func TestValidateExpiryYear(t *testing.T) {
	tests := []struct {
		name          string
		year          int
		expectedError bool
	}{
		{
			name:          "valid year 2024",
			year:          2024,
			expectedError: false,
		},
		{
			name:          "valid year 2030",
			year:          2030,
			expectedError: false,
		},
		{
			name:          "invalid year 0",
			year:          0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateExpiryYear(tt.year)

			if tt.expectedError && len(errors) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.expectedError && len(errors) > 0 {
				t.Errorf("expected no errors but got: %v", errors)
			}
			if tt.expectedError && len(errors) > 0 && errors[0].Field != "expiry_year" {
				t.Errorf("expected field 'expiry_year', got '%s'", errors[0].Field)
			}
		})
	}
}

func TestValidateExpiryDate(t *testing.T) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	tests := []struct {
		name           string
		month          int
		year           int
		expectedErrors int
		errorMessages  []string
	}{
		{
			name:           "valid future date",
			month:          12,
			year:           currentYear + 1,
			expectedErrors: 0,
		},
		{
			name:           "invalid past year",
			month:          6,
			year:           currentYear - 1,
			expectedErrors: 1,
			errorMessages:  []string{"expiry date must be in the future"},
		},
		{
			name:           "invalid current year past month",
			month:          currentMonth - 1,
			year:           currentYear,
			expectedErrors: 1,
			errorMessages:  []string{"expiry date must be in the future"},
		},
		{
			name:           "invalid month 0",
			month:          0,
			year:           currentYear + 1,
			expectedErrors: 1,
			errorMessages:  []string{"expiry month must be between 1-12"},
		},
		{
			name:           "invalid month 13",
			month:          13,
			year:           currentYear + 1,
			expectedErrors: 1,
			errorMessages:  []string{"expiry month must be between 1-12"},
		},
		{
			name:           "invalid year 0",
			month:          6,
			year:           0,
			expectedErrors: 1,
			errorMessages:  []string{"expiry year is required"},
		},
		{
			name:           "invalid month and year",
			month:          0,
			year:           0,
			expectedErrors: 2,
			errorMessages:  []string{"expiry month must be between 1-12", "expiry year is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateExpiryDate(tt.month, tt.year)

			if len(errors) != tt.expectedErrors {
				t.Errorf("expected %d errors, got %d: %v", tt.expectedErrors, len(errors), errors)
			}

			for i, msg := range tt.errorMessages {
				if i < len(errors) && errors[i].Message != msg {
					t.Errorf("expected error message '%s' at index %d, got '%s'", msg, i, errors[i].Message)
				}
			}
		})
	}
}

func TestValidateCurrency(t *testing.T) {
	tests := []struct {
		name          string
		currency      string
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "valid USD",
			currency:      "USD",
			expectedError: false,
		},
		{
			name:          "valid EUR",
			currency:      "EUR",
			expectedError: false,
		},
		{
			name:          "valid GBP",
			currency:      "GBP",
			expectedError: false,
		},
		{
			name:          "empty currency",
			currency:      "",
			expectedError: true,
			errorMessage:  "currency is required",
		},
		{
			name:          "invalid currency BRL",
			currency:      "BRL",
			expectedError: true,
			errorMessage:  "currency must be one of: USD, GBP, EUR",
		},
		{
			name:          "invalid lowercase usd",
			currency:      "usd",
			expectedError: true,
			errorMessage:  "currency must be one of: USD, GBP, EUR",
		},
		{
			name:          "invalid currency JPY",
			currency:      "JPY",
			expectedError: true,
			errorMessage:  "currency must be one of: USD, GBP, EUR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateCurrency(tt.currency)

			if tt.expectedError {
				if len(errors) == 0 {
					t.Errorf("expected error but got none")
				} else {
					if errors[0].Field != "currency" {
						t.Errorf("expected field 'currency', got '%s'", errors[0].Field)
					}
					if errors[0].Message != tt.errorMessage {
						t.Errorf("expected message '%s', got '%s'", tt.errorMessage, errors[0].Message)
					}
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("expected no errors but got: %v", errors)
				}
			}
		})
	}
}

func TestValidateCvv(t *testing.T) {
	tests := []struct {
		name           string
		cvv            string
		expectedErrors int
		errorMessages  []string
	}{
		{
			name:           "valid 3 digit cvv",
			cvv:            "123",
			expectedErrors: 0,
		},
		{
			name:           "valid 4 digit cvv",
			cvv:            "1234",
			expectedErrors: 0,
		},
		{
			name:           "empty cvv",
			cvv:            "",
			expectedErrors: 1,
			errorMessages:  []string{"cvv is required"},
		},
		{
			name:           "cvv too short",
			cvv:            "12",
			expectedErrors: 1,
			errorMessages:  []string{"cvv must be 3-4 characters long"},
		},
		{
			name:           "cvv too long",
			cvv:            "12345",
			expectedErrors: 1,
			errorMessages:  []string{"cvv must be 3-4 characters long"},
		},
		{
			name:           "cvv with letters",
			cvv:            "12a",
			expectedErrors: 1,
			errorMessages:  []string{"cvv must contain only numeric characters"},
		},
		{
			name:           "cvv with special characters",
			cvv:            "12@",
			expectedErrors: 1,
			errorMessages:  []string{"cvv must contain only numeric characters"},
		},
		{
			name:           "cvv with letters and wrong length",
			cvv:            "ab",
			expectedErrors: 2,
			errorMessages:  []string{"cvv must be 3-4 characters long", "cvv must contain only numeric characters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateCvv(tt.cvv)

			if len(errors) != tt.expectedErrors {
				t.Errorf("expected %d errors, got %d: %v", tt.expectedErrors, len(errors), errors)
			}

			for i, msg := range tt.errorMessages {
				if i < len(errors) {
					if errors[i].Field != "cvv" {
						t.Errorf("expected field 'cvv', got '%s'", errors[i].Field)
					}
					if errors[i].Message != msg {
						t.Errorf("expected message '%s', got '%s'", msg, errors[i].Message)
					}
				}
			}
		})
	}
}

func TestValidatePaymentRequest(t *testing.T) {
	now := time.Now()
	currentYear := now.Year()

	tests := []struct {
		name           string
		request        dto.PaymentRequest
		expectedErrors int
	}{
		{
			name: "valid payment request",
			request: dto.PaymentRequest{
				CardNumber:  "4532015112830366",
				ExpiryMonth: 12,
				ExpiryYear:  currentYear + 1,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			},
			expectedErrors: 0,
		},
		{
			name: "invalid payment request - all fields invalid",
			request: dto.PaymentRequest{
				CardNumber:  "",
				ExpiryMonth: 0,
				ExpiryYear:  0,
				Currency:    "",
				Amount:      0,
				Cvv:         "",
			},
			expectedErrors: 6,
		},
		{
			name: "invalid payment request - invalid card number",
			request: dto.PaymentRequest{
				CardNumber:  "123",
				ExpiryMonth: 12,
				ExpiryYear:  currentYear + 1,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			},
			expectedErrors: 1,
		},
		{
			name: "invalid payment request - invalid currency",
			request: dto.PaymentRequest{
				CardNumber:  "4532015112830366",
				ExpiryMonth: 12,
				ExpiryYear:  currentYear + 1,
				Currency:    "BRL",
				Amount:      1000,
				Cvv:         "123",
			},
			expectedErrors: 1,
		},
		{
			name: "invalid payment request - invalid amount",
			request: dto.PaymentRequest{
				CardNumber:  "4532015112830366",
				ExpiryMonth: 12,
				ExpiryYear:  currentYear + 1,
				Currency:    "USD",
				Amount:      -100,
				Cvv:         "123",
			},
			expectedErrors: 1,
		},
		{
			name: "invalid payment request - invalid cvv",
			request: dto.PaymentRequest{
				CardNumber:  "4532015112830366",
				ExpiryMonth: 12,
				ExpiryYear:  currentYear + 1,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "12",
			},
			expectedErrors: 1,
		},
		{
			name: "invalid payment request - expired date",
			request: dto.PaymentRequest{
				CardNumber:  "4532015112830366",
				ExpiryMonth: 1,
				ExpiryYear:  currentYear - 1,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			},
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidatePaymentRequest(tt.request)

			if len(errors) != tt.expectedErrors {
				t.Errorf("expected %d errors, got %d: %v", tt.expectedErrors, len(errors), errors)
			}
		})
	}
}

func TestMergeValidationErrors(t *testing.T) {
	tests := []struct {
		name           string
		errorSlices    [][]dto.ValidationError
		expectedLength int
	}{
		{
			name:           "no errors",
			errorSlices:    [][]dto.ValidationError{},
			expectedLength: 0,
		},
		{
			name: "single error slice",
			errorSlices: [][]dto.ValidationError{
				{
					{Field: "field1", Message: "error1"},
				},
			},
			expectedLength: 1,
		},
		{
			name: "multiple error slices",
			errorSlices: [][]dto.ValidationError{
				{
					{Field: "field1", Message: "error1"},
				},
				{
					{Field: "field2", Message: "error2"},
				},
			},
			expectedLength: 2,
		},
		{
			name: "empty and non-empty slices",
			errorSlices: [][]dto.ValidationError{
				{},
				{
					{Field: "field1", Message: "error1"},
				},
				{},
			},
			expectedLength: 1,
		},
		{
			name: "all empty slices",
			errorSlices: [][]dto.ValidationError{
				{},
				{},
				{},
			},
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeValidationErrors(tt.errorSlices...)

			if len(result) != tt.expectedLength {
				t.Errorf("expected %d errors, got %d", tt.expectedLength, len(result))
			}
		})
	}
}
