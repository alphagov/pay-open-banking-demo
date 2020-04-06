package truelayer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Config struct {
	AuthURL           string
	PayURL            string
	ClientID          string
	ClientSecret      string
	BankAccountNumber string
	BankSortCode      string
}

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
	RemitterProviderID           string `json:"remitter_provider_id"`
	DirectBankLink               bool   `json:"direct_bank_link"`
}

type TrueLayer struct {
	config *Config
	token  AccessTokenResponse
}

func NewTruelayer(config *Config) (*TrueLayer, error) {
	token, err := generatePaymentToken(config)
	if err != nil {
		return nil, err
	}

	return &TrueLayer{config: config, token: token}, nil
}

// GeneratePaymentToken Generates a token and stores it in the environment variables
func generatePaymentToken(config *Config) (AccessTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("scope", "payments")
	data.Set("grant_type", "client_credentials")

	accessTokenResponse := AccessTokenResponse{}
	req, err := http.NewRequest("POST", config.AuthURL+"/connect/token", strings.NewReader(data.Encode()))
	if err != nil {
		return accessTokenResponse, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return accessTokenResponse, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		return accessTokenResponse, err
	}

	log.Printf("Got TrueLayer token, expires in %d", accessTokenResponse.ExpiresIn)
	return accessTokenResponse, err
}

// CreateSinglePayment creates a payment in truelayer
func (trueLayer *TrueLayer) CreateSinglePayment(request *SinglePaymentRequest) (SinglePaymentResponse, error) {
	request.BeneficiarySortCode = trueLayer.config.BankSortCode
	request.BeneficiaryAccountNumber = trueLayer.config.BankAccountNumber
	
	paymentResponse := SinglePaymentResponse{}
	marshalled, err := json.Marshal(request)
	if err != nil {
		return paymentResponse, err
	}

	req, err := http.NewRequest("POST", trueLayer.config.PayURL+"/single-immediate-payments", bytes.NewBuffer(marshalled))
	req.Header.Set("Authorization", "Bearer "+trueLayer.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return paymentResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var data map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(body, &data)
		if err != nil {
			return paymentResponse, err
		}
		log.Printf("TrueLayer returned status code: %d", resp.StatusCode)
		log.Print(data)
		err = errors.New("Failed to create payment with TrueLayer")
		return paymentResponse, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return paymentResponse, err
	}
	return paymentResponse, err
}

// GetSinglePaymentInfo Gets the information for a single payment by the simpID
func (trueLayer *TrueLayer) GetSinglePaymentInfo(simpID string) (SinglePaymentResponse, error) {
	paymentResponse := SinglePaymentResponse{}
	req, err := http.NewRequest("GET", trueLayer.config.PayURL+"/single-immediate-payments/"+simpID, nil)
	if err != nil {
		return paymentResponse, err
	}

	req.Header.Set("Authorization", "Bearer "+trueLayer.token.AccessToken)
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

// GetSinglePaymentStatuses Gets the statuses for a single simpID and returns them
func (trueLayer *TrueLayer) GetSinglePaymentStatuses(simpID string) (PaymentStatusesResponse, error) {
	paymentStatuses := PaymentStatusesResponse{}

	req, err := http.NewRequest("GET", trueLayer.config.PayURL+"/single-immediate-payments/"+simpID+"/statuses", nil)
	if err != nil {
		return paymentStatuses, err
	}

	req.Header.Set("Authorization", "Bearer "+trueLayer.token.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return paymentStatuses, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &paymentStatuses)
	if err != nil {
		return paymentStatuses, err
	}
	return paymentStatuses, err
}

// GetProviders requests a list of providers from TrueLayer that can currently proccess a SingleImmediatePayment
func (trueLayer *TrueLayer) GetProviders() ProvidersResponse {
	req, err := http.NewRequest("GET", trueLayer.config.PayURL+"/providers?capability=SingleImmediatePayment", nil)
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
