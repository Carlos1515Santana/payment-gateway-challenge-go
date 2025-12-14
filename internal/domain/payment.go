package domain

type PaymentStatus string

const (
	PaymentStatusAuthorized PaymentStatus = "Authorized"
	PaymentStatusDeclined   PaymentStatus = "Declined"
	PaymentStatusRejected   PaymentStatus = "Rejected"
)

type Payment struct {
	ID       string
	Card     Card
	Currency string
	Amount   int
	Status   PaymentStatus
}

type Card struct {
	Number      string
	ExpiryMonth int
	ExpiryYear  int
}

func NewPayment(paymentID string, card Card, currency string, amount int, status PaymentStatus) *Payment {
	p := &Payment{
		ID:       paymentID,
		Card:     card,
		Currency: currency,
		Amount:   amount,
		Status:   status,
	}
	return p
}

func (c *Card) GetLastFourDigits() string {
	if len(c.Number) < 4 {
		return c.Number
	}
	return c.Number[len(c.Number)-4:]
}
