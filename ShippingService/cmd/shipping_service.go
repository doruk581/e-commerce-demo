package main

import (
	"log"
	"net/http"

	"github.com/ajanthan/go-ecommerce-demo/ShippingService/pkg/api"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	shippingEndpoint := new(api.ShippingEndpoint)
	router.HandleFunc("/shipping/getqote", shippingEndpoint.GetShippingQuote).Methods("POST")
	router.HandleFunc("/shipping/order", shippingEndpoint.ProcessShippingOrder).Methods("POST")
	log.Fatal(http.ListenAndServe(":8811", router))
}
