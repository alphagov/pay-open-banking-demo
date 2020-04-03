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

// AccessTokenResponse is the response we receive from TrueLayer when we request an access token
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// PaymentStatusesResponse is the wrapper object TrueLayer send to use when we request payment statuses
type PaymentStatusesResponse struct {
	Results []PaymentStatus `json:"results"`
	Status  string          `json:"status"`
}

// PaymentStatus is the object we receive from TrueLayer when we request payment statuses
type PaymentStatus struct {
	Status string `json:"status"`
	Date   string `json:"date"`
}

// SinglePaymentResponse is a wrapper for the results we receive when we receive information about a payment
type SinglePaymentResponse struct {
	PaymentResult []PaymentResult `json:"results"`
}

// ProvidersResponse is the wrapper object we receive from TrueLayer to list banks that are currently available
type ProvidersResponse struct {
	Results []Provider `json:"results"`
}

// Provider is the representation of a bank we receive from TrueLayer
type Provider struct {
	ID              string `json:"id"`
	Logo            string `json:"logo"`
	Icon            string `json:"icon"`
	DisplayableName string `json:"displayable_name"`
	MainBgColor     string `json:"main_bg_color"`
}

// PaymentResult is the representation of a payment we receive from TrueLayer
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

type SinglePaymentRequest struct {
	Amount                       int    `json:"amount"`
	Currency                     string `json:"currency"`
	BeneficiaryName              string `json:"beneficiary_name"`
	BeneficiaryReference         string `json:"beneficiary_reference"`
	BeneficiarySortCode          string `json:"beneficiary_sort_code"`
	BeneficiaryAccountNumber     string `json:"beneficiary_account_number"`
	BeneficiaryRemitterReference string `json:"remitter_reference"`
	RedirectURL                  string `json:"redirect_uri"`
}

// CreateSinglePayment creates a payment in truelayer
func CreateSinglePayment(request SinglePaymentRequest, accessToken string) (SinglePaymentResponse, error) {
	paymentResponse := SinglePaymentResponse{}
	marshalled, err := json.Marshal(request)
	if err != nil {
		return paymentResponse, err
	}

	baseTrueLayerPayURL := os.Getenv("TRUELAYER_PAY_URL")
	req, err := http.NewRequest("POST", baseTrueLayerPayURL+"/single-immediate-payments", bytes.NewBuffer(marshalled))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return paymentResponse, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return paymentResponse, err
	}
	return paymentResponse, err
}

// GetSinglePaymentInfo Gets the information for a single payment by the simpID
func GetSinglePaymentInfo(simpID string, accessToken string) SinglePaymentResponse {
	baseTrueLayerPayURL := os.Getenv("TRUELAYER_PAY_URL")
	req, err := http.NewRequest("GET", baseTrueLayerPayURL+"/single-immediate-payments/"+simpID, nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	paymentResponse := SinglePaymentResponse{}
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		panic(err)
	}
	return paymentResponse
}

// GetSinglePaymentStatuses Gets the statuses for a single simpID and returns them
func GetSinglePaymentStatuses(simpID string, accessToken string) PaymentStatusesResponse {
	baseTrueLayerPayURL := os.Getenv("TRUELAYER_PAY_URL")
	req, err := http.NewRequest("GET", baseTrueLayerPayURL+"/single-immediate-payments/"+simpID+"/statuses", nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	paymentStatuses := PaymentStatusesResponse{}
	err = json.Unmarshal(body, &paymentStatuses)
	if err != nil {
		panic(err)
	}
	return paymentStatuses
}

// GeneratePaymentToken Generates a token and stores it in the environment variables
func GeneratePaymentToken() AccessTokenResponse {
	baseTrueLayerAuthURL := os.Getenv("TRUELAYER_AUTH_URL")

	clientID := os.Getenv("TRUELAYER_CLIENT_ID")
	clientSecret := os.Getenv("TRUELAYER_CLIENT_SECRET")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "payments")
	data.Set("grant_type", "client_credentials")

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
	accessTokenResponse := AccessTokenResponse{}
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		panic(err)
	}

	return accessTokenResponse
}

// GetProviders requests a list of providers from TrueLayer that can currently proccess a SingleImmediatePayment
func GetProviders() ProvidersResponse {
	baseTrueLayerPayURL := os.Getenv("TRUELAYER_PAY_URL")

	req, err := http.NewRequest("GET", baseTrueLayerPayURL+"/providers?capability=SingleImmediatePayment", nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	providers := ProvidersResponse{}
	err = json.Unmarshal(body, &providers)
	if err != nil {
		panic(err)
	}
	return providers
}
