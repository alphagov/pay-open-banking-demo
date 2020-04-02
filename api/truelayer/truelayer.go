package truelayer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type PaymentStatusesResponse struct {
	Results []PaymentStatus `json:"results"`
	Status  string          `json:"status"`
}

type PaymentStatus struct {
	Status string `json:"status"`
	Date   string `json:"date"`
}

type SinglePaymentResponse struct {
	PaymentResult []PaymentResult `json:"results"`
	Status        string          `json:"status"`
}

type PaymentResult struct {
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

// CreateSinglePayment creates a payment in truelayer
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

	var baseTrueLayerPayURL = os.Getenv("PAY_OPEN_BANKING_TRUELAYER_PAY_URL")
	req, err := http.NewRequest("POST", baseTrueLayerPayURL+"/single-immediate-payments", bytes.NewBuffer(marshalled))
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

// GetSinglePaymentInfo Gets the information for a single payment by the simpID
func GetSinglePaymentInfo(simpID string) SinglePaymentResponse {
	var baseTrueLayerPayURL = os.Getenv("PAY_OPEN_BANKING_TRUELAYER_PAY_URL")
	req, err := http.NewRequest("GET", baseTrueLayerPayURL+"/single-immediate-payments/"+simpID, nil)

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

// GetSinglePaymentStatuses Gets the statuses for a single simpID and returns them
func GetSinglePaymentStatuses(simpID string) PaymentStatusesResponse {
	var baseTrueLayerPayURL = os.Getenv("PAY_OPEN_BANKING_TRUELAYER_PAY_URL")
	req, err := http.NewRequest("GET", baseTrueLayerPayURL+"/single-immediate-payments/"+simpID+"/statuses", nil)

	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAY_OPEN_BANKING_DEMO_TRUELAYER_TOKEN"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var paymentStatuses = PaymentStatusesResponse{}
	err = json.Unmarshal(body, &paymentStatuses)
	if err != nil {
		panic(err)
	}
	return paymentStatuses
}

// GeneratePaymentToken Generates a token and stores it in the environment variables
func GeneratePaymentToken() AccessTokenResponse {
	var baseTrueLayerAuthURL = os.Getenv("PAY_OPEN_BANKING_TRUELAYER_AUTH_URL")

	var clientID = os.Getenv("PAY_OPEN_BANKING_DEMO_CLIENT_ID")
	var clientSecret = os.Getenv("PAY_OPEN_BANKING_DEMO_CLIENT_SECRET")
	var scope = os.Getenv("PAY_OPEN_BANKING_DEMO_SCOPE")
	var grantType = os.Getenv("PAY_OPEN_BANKING_DEMO_GRANT_TYPE")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", scope)
	data.Set("grant_type", grantType)

	req, err := http.NewRequest("POST", baseTrueLayerAuthURL+"/connect/token", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var accessTokenResponse = AccessTokenResponse{}
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		panic(err)
	}
	os.Setenv("PAY_OPEN_BANKING_DEMO_TRUELAYER_TOKEN", accessTokenResponse.AccessToken)
	return accessTokenResponse
}
