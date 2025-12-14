package validators

import (
	"regexp"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/dto"
)

var (
	supportedCurrencies = map[string]bool{
		"USD": true,
		"EUR": true,
		"GBP": true,
	}

	numericRegex = regexp.MustCompile(`^[0-9]+$`)
)

func ValidatePaymentRequest(req dto.PaymentRequest) []dto.ValidationError {
	return mergeValidationErrors(
		ValidateCardNumber(req.CardNumber),
		ValidateExpiryDate(req.ExpiryMonth, req.ExpiryYear),
		ValidateCurrency(req.Currency),
		ValidateAmount(req.Amount),
		ValidateCvv(req.Cvv),
	)
}

func ValidateAmount(amount int) []dto.ValidationError {
	var errors []dto.ValidationError

	if amount <= 0 {
		errors = append(errors, dto.ValidationError{
			Field:   "amount",
			Message: "amount must be a positive integer",
		})
	}

	return errors
}

func ValidateCardNumber(cardNumber string) []dto.ValidationError {
	var errors []dto.ValidationError

	if cardNumber == "" {
		errors = append(errors, dto.ValidationError{
			Field:   "card_number",
			Message: "card number is required",
		})
	} else {
		if !numericRegex.MatchString(cardNumber) {
			errors = append(errors, dto.ValidationError{
				Field:   "card_number",
				Message: "card number must contain only numeric characters",
			})
		}
		if len(cardNumber) < 14 || len(cardNumber) > 19 {
			errors = append(errors, dto.ValidationError{
				Field:   "card_number",
				Message: "card number must be between 14-19 characters",
			})
		}
	}
	return errors
}

func ValidateExpiryDate(month int, year int) []dto.ValidationError {
	var errors []dto.ValidationError
	errors = append(errors, ValidateExpiryMonth(month)...)
	errors = append(errors, ValidateExpiryYear(year)...)

	if len(errors) > 0 {
		return errors
	}

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	if year < currentYear || (year == currentYear && month < currentMonth) {
		errors = append(errors, dto.ValidationError{
			Field:   "expiry_date",
			Message: "expiry date must be in the future",
		})
	}

	return errors
}
func ValidateExpiryMonth(month int) []dto.ValidationError {
	var errors []dto.ValidationError
	if month < 1 || month > 12 {
		errors = append(errors, dto.ValidationError{
			Field:   "expiry_month",
			Message: "expiry month must be between 1-12",
		})
	}

	return errors
}

func ValidateExpiryYear(year int) []dto.ValidationError {
	var errors []dto.ValidationError

	if year == 0 {
		errors = append(errors, dto.ValidationError{
			Field:   "expiry_year",
			Message: "expiry year is required",
		})
	}

	return errors
}

func ValidateCurrency(currency string) []dto.ValidationError {
	var errors []dto.ValidationError
	if currency == "" {
		errors = append(errors, dto.ValidationError{
			Field:   "currency",
			Message: "currency is required",
		})
		return errors
	}

	if !supportedCurrencies[currency] {
		errors = append(errors, dto.ValidationError{
			Field:   "currency",
			Message: "currency must be one of: USD, GBP, EUR",
		})
	}

	return errors
}

func ValidateCvv(cvv string) []dto.ValidationError {
	var errors []dto.ValidationError
	if cvv == "" {
		errors = append(errors, dto.ValidationError{
			Field:   "cvv",
			Message: "cvv is required",
		})
	} else {
		if len(cvv) < 3 || len(cvv) > 4 {
			errors = append(errors, dto.ValidationError{
				Field:   "cvv",
				Message: "cvv must be 3-4 characters long",
			})
		}
		if !numericRegex.MatchString(cvv) {
			errors = append(errors, dto.ValidationError{
				Field:   "cvv",
				Message: "cvv must contain only numeric characters",
			})
		}
	}
	return errors
}

func mergeValidationErrors(slicesOfErrs ...[]dto.ValidationError) []dto.ValidationError {
	var result []dto.ValidationError
	for _, errors := range slicesOfErrs {
		if len(errors) > 0 {
			result = append(result, errors...)
		}
	}
	return result
}
