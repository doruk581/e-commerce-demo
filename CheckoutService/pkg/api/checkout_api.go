package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	cart "github.com/ajanthan/go-ecommerce-demo/CartService/pkg/model"
	"github.com/ajanthan/go-ecommerce-demo/CheckoutService/pkg/model"
	payment "github.com/ajanthan/go-ecommerce-demo/PaymentService/pkg/model"
	product "github.com/ajanthan/go-ecommerce-demo/ProductCatalogService/pkg/data"
	shipping "github.com/ajanthan/go-ecommerce-demo/ShippingService/pkg/model"
)

type CheckoutEndpoint struct {
}

func (c *CheckoutEndpoint) Checkout(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, err)
		return
	}
	defer req.Body.Close()

	checkoutRequest := new(model.Order)
	unmarshalErr := json.Unmarshal(body, &checkoutRequest)

	if unmarshalErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, unmarshalErr)
		return
	}
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(randomSource)
	orderID := randomGen.Intn(1000000)
	orderResult := model.OrderResult{}
	orderResult.OrderID = fmt.Sprint(orderID)
	//Get the cart
	cart, cartErr := getCart(checkoutRequest.UserID)
	if cartErr != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, cartErr)
		return
	}
	//Get total cost estimation
	cost, costCalErr := calculateCost(cart)
	if costCalErr != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, costCalErr)
		return
	}
	// Get shipping cost
	shippingCost, shippingCostCalErr := getShippingCostEstimate(checkoutRequest.Address, cart)
	if shippingCostCalErr != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, shippingCostCalErr)
		return
	}
	cost = +shippingCost
	//Charge cc
	_, paymentErr := charge(checkoutRequest.CreditCardInfo, cost)
	if paymentErr != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, paymentErr)
		return
	}
	//Place shipping order
	trackingID, shippingOrderError := placeOrder(checkoutRequest.Address, cart)
	if shippingOrderError != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, shippingOrderError)
		return
	}
	//Empty cart
	emptyCartErr := emptyCart(checkoutRequest.UserID)

	if emptyCartErr != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, emptyCartErr)
		return
	}
	//send the result
	orderResult.Address = checkoutRequest.Address
	orderResult.Cart = cart
	orderResult.Cost = float64(cost)
	orderResult.TrackingID = trackingID
	respBytes, respMarshalErr := json.Marshal(orderResult)

	if respMarshalErr != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, respMarshalErr)
		return
	}
	res.Write(respBytes)
	res.Header().Set("Content-Type", "application/json")

}

func getCart(userID string) (cart.Cart, error) {
	cart := cart.Cart{}
	resp, CartErr := http.Get("http://localhost:8889/cart/" + userID)
	//handle 500
	if CartErr != nil {
		return cart, CartErr
	}
	respBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return cart, readErr
	}
	unmarshalErr := json.Unmarshal(respBytes, &cart)
	if unmarshalErr != nil {
		return cart, unmarshalErr
	}
	return cart, nil
}

func calculateCost(c cart.Cart) (float32, error) {
	var cost float32
	for _, item := range c.Items {
		resp, catelogErr := http.Get("http://localhost:8888/product/" + item.ProductID)
		if catelogErr != nil {
			return cost, catelogErr
		}
		respBytes, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return cost, readErr
		}
		defer resp.Body.Close()
		product := product.Product{}
		unMarshalErr := json.Unmarshal(respBytes, &product)
		if unMarshalErr != nil {
			return cost, unMarshalErr
		}
		cost = cost + product.Price*float32(item.Quantity)
	}
	return cost, nil
}

func getShippingCostEstimate(address shipping.Address, cart cart.Cart) (float32, error) {
	var shipingCost float32
	shippingQuoteRequest := shipping.ShippingQuoteRequest{}
	shippingQuoteRequest.Address = address
	shippingQuoteRequest.Cart = cart
	bodyBytes, marshalErr := json.Marshal(shippingQuoteRequest)
	if marshalErr != nil {
		return shipingCost, marshalErr
	}
	reader := bytes.NewReader(bodyBytes)
	resp, nwErr := http.Post("http://localhost:8811/shipping/getqote", "application/json", reader)
	if nwErr != nil {
		return shipingCost, nwErr
	}
	shippingQuoteResponse := shipping.ShippingQuoteResponse{}
	respBytes, resErr := ioutil.ReadAll(resp.Body)
	if resErr != nil {
		return shipingCost, resErr
	}
	unMarshalErr := json.Unmarshal(respBytes, &shippingQuoteResponse)

	if unMarshalErr != nil {
		return shipingCost, unMarshalErr
	}
	shipingCost = float32(shippingQuoteResponse.Cost)
	return shipingCost, nil
}
func charge(creditCard payment.CreditCard, amount float32) (string, error) {
	var transactionID string
	paymentRequest := payment.PaymentRequest{}
	paymentRequest.CreditCardInfo = creditCard
	paymentRequest.Amount = float64(amount)
	bodyBytes, marshalErr := json.Marshal(paymentRequest)
	if marshalErr != nil {
		return transactionID, marshalErr
	}
	reader := bytes.NewReader(bodyBytes)
	resp, nwErr := http.Post("http://localhost:8810/payment", "application/json", reader)
	if nwErr != nil {
		return transactionID, nwErr
	}
	paymentResponse := payment.PaymentReponse{}
	bodyBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return transactionID, readErr
	}
	unMarshalErr := json.Unmarshal(bodyBytes, &paymentResponse)
	if unMarshalErr != nil {
		return transactionID, unMarshalErr
	}
	transactionID = paymentResponse.TransactionID
	return transactionID, nil
}

func placeOrder(address shipping.Address, cart cart.Cart) (string, error) {
	var trackingID string
	shippingOrderRequest := shipping.ShippingOrderRequest{}
	shippingOrderRequest.Cart = cart
	shippingOrderRequest.Address = address
	bodyBytes, marshalErr := json.Marshal(shippingOrderRequest)
	if marshalErr != nil {
		return trackingID, marshalErr
	}
	reader := bytes.NewReader(bodyBytes)
	resp, nwErr := http.Post("http://localhost:8811/shipping/order", "application/json", reader)
	if nwErr != nil {
		return trackingID, nwErr
	}
	shippingOrderResponse := shipping.ShippingOrderResponse{}
	bodyBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return trackingID, readErr
	}
	unMarshalErr := json.Unmarshal(bodyBytes, &shippingOrderResponse)
	if unMarshalErr != nil {
		return trackingID, unMarshalErr
	}
	trackingID = shippingOrderResponse.TrackingID
	return trackingID, nil
}

func emptyCart(userID string) error {
	req, rqErr := http.NewRequest(http.MethodDelete, "http://localhost:8889/cart/"+userID, nil)
	if rqErr != nil {
		return rqErr
	}
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return respErr
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.New("Error in emptying cart")
	}
	return nil
}
