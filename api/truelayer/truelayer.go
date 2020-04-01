package truelayer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type SinglePaymentResponse struct {
	PaymentResult []paymentResult `json:"results"`
	Status        string          `json:"status"`
}

type paymentResult struct {
	SimpID                   string `json:"simp_id"`
	AuthURI                  string `json:"auth_uri"`
	CreatedAt                string `json:"created_at"`
	Amount                   int    `json:"amount"`
	Currency                 string `json:"currency"`
	BeneficiaryReference     string `json:"beneficiary_reference"`
	BeneficiaryName          string `json:"beneficiary_name"`
	BeneficiarySortCode      string `json:"beneficiary_sort_code"`
	BeneficiaryAccountNumber string `json:"beneficiary_account_number"`
	RemitterReference        string `json:"remitter_reference"`
	RedirectURI              string `json:"redirect_uri"`
	Status                   string `json:"status"`
}

type singlePaymentRequest struct {
	Amount                       int    `json:"amount"`
	Currency                     string `json:"currency"`
	BeneficiaryName              string `json:"beneficiary_name"`
	BeneficiaryReference         string `json:"beneficiary_reference"`
	BeneficiarySortCode          string `json:"beneficiary_sort_code"`
	BeneficiaryAccountNumber     string `json:"beneficiary_account_number"`
	BeneficiaryRemitterReference string `json:"remitter_reference"`
	RedirectURL                  string `json:"redirect_uri"`
}

func newSinglePaymentRequest(amount int,
	currency string,
	beneficiaryName string,
	beneficiaryReference string,
	beneficiarySortCode string,
	beneficiaryAccountNumber string,
	beneficiaryRemitterReference string,
	redirectURL string) *singlePaymentRequest {
	return &singlePaymentRequest{Amount: amount,
		Currency:                     currency,
		BeneficiaryName:              beneficiaryName,
		BeneficiaryReference:         beneficiaryReference,
		BeneficiarySortCode:          beneficiarySortCode,
		BeneficiaryAccountNumber:     beneficiaryAccountNumber,
		BeneficiaryRemitterReference: beneficiaryRemitterReference,
		RedirectURL:                  redirectURL,
	}
}

var baseTrueLayerURL = os.Getenv("PAY_OPEN_BANKING_TRUELAYER_URL")

func CreateSinglePayment(amount int,
	currency string,
	beneficiaryName string,
	beneficiaryReference string,
	beneficiarySortCode string,
	beneficiaryAccountNumber string,
	beneficiaryRemitterReference string,
	redirectURL string) SinglePaymentResponse {

	var jsonBody = newSinglePaymentRequest(amount,
		currency,
		beneficiaryName,
		beneficiaryReference,
		beneficiarySortCode,
		beneficiarySortCode,
		beneficiaryRemitterReference,
		redirectURL)

	marshalled, err := json.Marshal(jsonBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", baseTrueLayerURL+"/single-immediate-payments", bytes.NewBuffer(marshalled))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAY_OPEN_BANKING_DEMO_TRUELAYER_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var paymentResponse = SinglePaymentResponse{}
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		panic(err)
	}
	return paymentResponse
}

func GetSinglePaymentInfo(simpID string) SinglePaymentResponse {
	req, err := http.NewRequest("GET", baseTrueLayerURL+"/single-immediate-payments/"+simpID, nil)

	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAY_OPEN_BANKING_DEMO_TRUELAYER_TOKEN"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var paymentResponse = SinglePaymentResponse{}
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		panic(err)
	}
	return paymentResponse
}
