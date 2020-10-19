package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	product "github.com/ajanthan/go-ecommerce-demo/ProductCatalogService/pkg/data"
)

//ProductAPI is data structure hold the static data
type ProductAPI struct {
	CatelogMap map[string]product.Product
}

func (p *ProductAPI) GetProductHandler(res http.ResponseWriter, req *http.Request) {
	args := mux.Vars(req)
	productID := args["id"]
	product, isExist := p.CatelogMap[productID]
	if isExist {
		data, jsonerr := json.Marshal(product)
		if jsonerr != nil {
			log.Fatal(jsonerr)
		}
		_, writeErr := res.Write(data)
		if writeErr != nil {
			log.Fatal(writeErr)
		}
	} else {
		res.WriteHeader(404)
	}
}

func (p *ProductAPI) GetAllProductsHandler(res http.ResponseWriter, req *http.Request) {

	var products []product.Product
	for _, value := range p.CatelogMap {
		products = append(products, value)
	}
	data, jsonerr := json.Marshal(products)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	_, writeErr := res.Write(data)
	if writeErr != nil {
		log.Fatal(writeErr)
	}
}

func (p *ProductAPI) SearchProductHandler(res http.ResponseWriter, req *http.Request) {

	queryString := req.URL.Query().Get("query")
	var matchingProducts []product.Product

	for id, prod := range p.CatelogMap {
		if strings.Contains(prod.Name, queryString) {
			matchingProducts = append(matchingProducts, p.CatelogMap[id])
		}
	}
	data, jsonerr := json.Marshal(matchingProducts)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	_, writeErr := res.Write(data)
	if writeErr != nil {
		log.Fatal(writeErr)
	}

}
