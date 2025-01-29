package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RecaptchaAPI struct {
	SecretKey string
}

func NewRecaptchaAPI(recaptchaSecretKey string) *RecaptchaAPI {
	return &RecaptchaAPI{
		SecretKey: recaptchaSecretKey,
	}
}

type RecaptchaV2Response struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"` // ISO timestamp when the challenge was completed
	Hostname    string   `json:"hostname"`     // The hostname where reCAPTCHA was executed
	ErrorCodes  []string `json:"error-codes"`  // Optional error codes
}

func (r *RecaptchaAPI) Verify(responseToken string) (bool, error) {
	// Prepare the API request to Google's reCAPTCHA endpoint
	apiURL := "https://www.google.com/recaptcha/api/siteverify"

	data := url.Values{}
	data.Set("secret", r.SecretKey)
	data.Set("response", responseToken)

	// Make the POST request to Google's reCAPTCHA API
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return false, fmt.Errorf("failed to send request to reCAPTCHA API: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read reCAPTCHA response: %w", err)
	}

	// Unmarshal the response JSON into the RecaptchaV2Response struct
	var recaptchaResponse RecaptchaV2Response
	err = json.Unmarshal(body, &recaptchaResponse)
	if err != nil {
		return false, fmt.Errorf("failed to parse reCAPTCHA response JSON: %w", err)
	}

	return recaptchaResponse.Success, nil
}
