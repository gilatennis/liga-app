package common

import (
	"encoding/json"
	"fmt"
)

type AlertType string

const (
	AlertTypeSuccess AlertType = "success"
	AlertTypeError   AlertType = "error"
)

type AlertPayload struct {
	Alert AlertPayloadDetail `json:"alert"`
}

type AlertPayloadDetail struct {
	Type    AlertType `json:"type"`
	Message string    `json:"message"`
}

func NewAlertPayload(alertType AlertType, message string) *AlertPayload {
	return &AlertPayload{
		Alert: AlertPayloadDetail{
			Type:    alertType,
			Message: message,
		},
	}
}

func (a *AlertPayload) ToJSON() (string, error) {
	msg, err := json.Marshal(a)
	if err != nil {
		return "", fmt.Errorf("error marshalling alert payload: %w", err)
	}
	return string(msg), nil
}

type AlertPayloadWithRedirect struct {
	Alert AlertPayloadWithRedirectDetail `json:"alert"`
}

type AlertPayloadWithRedirectDetail struct {
	Type     AlertType `json:"type"`
	Message  string    `json:"message"`
	Redirect string    `json:"redirect"`
}

func NewAlertPayloadWithRedirect(alertType AlertType, message string, redirect string) *AlertPayloadWithRedirect {
	return &AlertPayloadWithRedirect{
		Alert: AlertPayloadWithRedirectDetail{
			Type:     alertType,
			Message:  message,
			Redirect: redirect,
		},
	}
}

func (a *AlertPayloadWithRedirect) ToJSON() (string, error) {
	msg, err := json.Marshal(a)
	if err != nil {
		return "", fmt.Errorf("error marshalling alert payload: %w", err)
	}
	return string(msg), nil
}
