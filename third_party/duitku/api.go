package duitku

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"liga-app/config"
	"liga-app/interfaces/third_party/duitku"
	"net/http"
	"time"

	"github.com/k0kubun/pp"
	"github.com/moul/http2curl"
)

var logTag = "[third_party.duitku.api]"

type API struct {
	cfg *config.Config
}

func NewAPI(cfg *config.Config) *API {
	return &API{cfg: cfg}
}

func (a *API) CreateInvoice(ctx context.Context, payload *duitku.CreateInvoicePayload) (*duitku.CreateInvoiceResponse, error) {
	url := a.cfg.Duitku.ApiURL + "/api/merchant/createInvoice"
	apiKey := a.cfg.Duitku.ApiKey

	client := http.Client{}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("%s CreateInvoice(): failed to marshal payload: %w", logTag, err)
	}
	fmt.Println("DEBUG: payloadJson", string(payloadJson))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadJson))
	if err != nil {
		return nil, fmt.Errorf("%s CreateInvoice(): failed to create new request: %w", logTag, err)
	}

	timestamp := time.Now().UnixMilli()
	stringToSign := fmt.Sprintf("%s%d%s", a.cfg.Duitku.MerchantCode, timestamp, apiKey)
	h := sha256.New()
	h.Write([]byte(stringToSign))
	signature := h.Sum(nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-duitku-merchantcode", a.cfg.Duitku.MerchantCode)
	req.Header.Set("x-duitku-signature", hex.EncodeToString(signature))
	req.Header.Set("x-duitku-timestamp", fmt.Sprintf("%d", timestamp))

	curlCommand, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get curl command: %s", err)
	}

	fmt.Println("curlCommand", curlCommand)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s CreateInvoice(): failed to send request: %w", logTag, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s CreateInvoice(): failed to read response body: %w", logTag, err)
	}

	/*
		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, respBody, "", "\t")
		if error != nil {
			return nil, fmt.Errorf("JSON parse error: %w", error)
		}
		fmt.Println("prettyJSON", prettyJSON.String())
	*/
	fmt.Println("respBody", string(respBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s CreateInvoice(): unexpected status code: %d", logTag, resp.StatusCode)
	}

	var response duitku.CreateInvoiceResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("%s CreateInvoice(): failed to decode response: %w", logTag, err)
	}

	pp.Println("response", response)

	return &response, nil
}
