package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/ajanthan/go-ecommerce-demo/PaymentService/model"
)

type PaymentEndpoint struct {
}

func (p *PaymentEndpoint) Charge(res http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, err)
		return
	}
	defer req.Body.Close()

	paymentRequest := new(model.PaymentRequest)
	unmarshalErr := json.Unmarshal(body, &paymentRequest)

	if unmarshalErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, unmarshalErr)
		return
	}

	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(randomSource)
	transactionID := randomGen.Intn(1000000)
	log.Printf("Processing payment: transactionID: %d, %s", transactionID, *paymentRequest)
	paymentResponse := model.PaymentReponse{}
	paymentResponse.TransactionID = fmt.Sprint(transactionID)
	resBytes, marshelErr := json.Marshal(paymentResponse)
	if marshelErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, marshelErr)
		return
	}
	res.Write(resBytes)
	res.Header().Set("Content-Type", "application/json")

}
