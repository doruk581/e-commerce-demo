package model

import (
	"github.com/ajanthan/go-ecommerce-demo/CartService/model"
	payment "github.com/ajanthan/go-ecommerce-demo/PaymentService/model"
	shipping "github.com/ajanthan/go-ecommerce-demo/ShippingService/model"
)

type Order struct {
	UserID         string             `json:"userID"`
	Email          string             `json:"email"`
	Address        shipping.Address   `json:"address"`
	CreditCardInfo payment.CreditCard `json:"creditCard"`
}

type OrderResult struct {
	OrderID    string
	TrackingID string
	Address    shipping.Address
	Cart       model.Cart
	Cost       float64
}
