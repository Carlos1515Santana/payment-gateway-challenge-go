package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCard_GetLastFourDigits(t *testing.T) {
	tests := []struct {
		name           string
		cardNumber     string
		expectedResult string
	}{
		{
			name:           "standard 16-digit card",
			cardNumber:     "1234567890123456",
			expectedResult: "3456",
		},
		{
			name:           "14-digit card",
			cardNumber:     "12345678901234",
			expectedResult: "1234",
		},
		{
			name:           "19-digit card",
			cardNumber:     "1234567890123456789",
			expectedResult: "6789",
		},
		{
			name:           "exactly 4 digits",
			cardNumber:     "1234",
			expectedResult: "1234",
		},
		{
			name:           "less than 4 digits",
			cardNumber:     "123",
			expectedResult: "123",
		},
		{
			name:           "single digit",
			cardNumber:     "1",
			expectedResult: "1",
		},
		{
			name:           "empty card number",
			cardNumber:     "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := &Card{
				Number:      tt.cardNumber,
				ExpiryMonth: 12,
				ExpiryYear:  2025,
			}

			result := card.GetLastFourDigits()

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
