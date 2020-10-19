package main

import (
	"log"
	"net/http"

	"github.com/ajanthan/go-ecommerce-demo/ProductCatalogService/api"
	"github.com/gorilla/mux"

	"github.com/ajanthan/go-ecommerce-demo/ProductCatalogService/loader"
)

func main() {
	productAPI := &api.ProductAPI{}
	productAPI.CatelogMap = loader.LoadProductCatelog("products.json")
	router := mux.NewRouter()

	router.HandleFunc("/product/{id}", productAPI.GetProductHandler).Methods("GET")
	router.HandleFunc("/product", productAPI.GetAllProductsHandler).Methods("GET")
	router.HandleFunc("/search", productAPI.SearchProductHandler).Methods("GET").Queries("query", "{query}")
	log.Fatal(http.ListenAndServe(":8888", router))
}
