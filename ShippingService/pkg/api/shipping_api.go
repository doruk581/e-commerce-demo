package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/ajanthan/go-ecommerce-demo/ShippingService/pkg/model"
)

type ShippingEndpoint struct {
}

func (s *ShippingEndpoint) GetShippingQuote(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, err)
		return
	}
	defer req.Body.Close()

	shippingQuoteRequest := new(model.ShippingQuoteRequest)
	unmarshalErr := json.Unmarshal(body, &shippingQuoteRequest)

	if unmarshalErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, unmarshalErr)
		return
	}
	var cost float64
	count := len(shippingQuoteRequest.Cart.Items)
	if count != 0 {
		cost = math.Pow(3, (1 + 0.2*float64(count)))
	}
	shippingQuoteResponse := new(model.ShippingQuoteResponse)
	shippingQuoteResponse.Cost = cost
	response, marshalErr := json.Marshal(shippingQuoteResponse)
	if marshalErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, marshalErr)
		return
	}
	if _, writeErr := res.Write(response); writeErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, writeErr)
		return
	}

}

func (s *ShippingEndpoint) ProcessShippingOrder(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, err)
		return
	}
	defer req.Body.Close()

	shippingOrderRequest := new(model.ShippingOrderRequest)
	unmarshalErr := json.Unmarshal(body, &shippingOrderRequest)

	if unmarshalErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, unmarshalErr)
		return
	}

	shippingOrderResponse := new(model.ShippingOrderResponse)

	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(randomSource)
	transactionID := randomGen.Intn(1000000)
	log.Println("TransactionID :", transactionID)
	shippingOrderResponse.TrackingID = fmt.Sprint(transactionID)
	response, writeErr := json.Marshal(shippingOrderResponse)
	if writeErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, writeErr)
		return
	}
	log.Println("Processing shipping order: Address", shippingOrderRequest.Address, ",Order{", shippingOrderRequest.Cart, "}, TrackingID: ", transactionID)
	res.Write(response)
}
