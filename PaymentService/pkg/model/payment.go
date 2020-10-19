package model

import "fmt"

type CreditCard struct {
	Number          string
	CSV             int
	ExpirationYear  int
	ExpirationMonth int
}

func (c CreditCard) String() string {
	creditCardStr := "number:%s,csv: %d,expirationMonth: %d,expirationYear:%d"
	return fmt.Sprintf(creditCardStr, c.Number, c.CSV, c.ExpirationMonth, c.ExpirationYear)
}

type PaymentRequest struct {
	CreditCardInfo CreditCard `json:"creditCard"`
	Amount         float64    `json:"amount"`
}
type PaymentReponse struct {
	TransactionID string
}

func (p PaymentRequest) String() string {
	paymentRequestStr := "creditCard:{%s},amount:%f"
	return fmt.Sprintf(paymentRequestStr, p.CreditCardInfo, p.Amount)
}
