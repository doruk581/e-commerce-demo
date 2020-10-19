package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ajanthan/go-ecommerce-demo/CartService/pkg/data"
	"github.com/ajanthan/go-ecommerce-demo/CartService/pkg/model"
)

type CartAPI struct {
	Repository *data.CartRepository
}

func (c *CartAPI) AddCartHandler(res http.ResponseWriter, req *http.Request) {
	cart := new(model.Cart)
	cartByte, reqBodyErr := ioutil.ReadAll(req.Body)
	if reqBodyErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()
	jsonErr := json.Unmarshal(cartByte, &cart)
	if jsonErr != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, jsonErr)
		return
	}
	for _, item := range cart.Items {
		c.Repository.AddOrUpdateCartItem(cart.UserID, item)
	}
	res.WriteHeader(http.StatusAccepted)
}

func (c *CartAPI) GetCartHandler(res http.ResponseWriter, req *http.Request) {
	args := mux.Vars(req)
	userID := args["userID"]
	cart := c.Repository.GetCart(userID)

	if len(cart.Items) != 0 {
		cartByte, jsonErr := json.Marshal(cart)
		if jsonErr != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(cartByte)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

func (c *CartAPI) EmptyCartHandler(res http.ResponseWriter, req *http.Request) {
	args := mux.Vars(req)
	userID := args["userID"]
	c.Repository.EmptyCart(userID)
	res.WriteHeader(http.StatusAccepted)
}
