package mailtrap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"liga-app/config"
	"log"
	"net/http"
	"time"
)

type MailTrap struct {
	ApiURL string
	Token  string
}

func NewMailTrap(config *config.Config) *MailTrap {
	return &MailTrap{
		ApiURL: config.MailTrap.ApiURL,
		Token:  config.MailTrap.Token,
	}
}

func (m *MailTrap) SendEmail(ctx context.Context, from string, to string, subject string, content string) error {
	msg := fmt.Sprintf(
		`{"from": {"email": "%s", "name": "%s"},
    "to": [{"email": "%s"}],
    "subject": "%s",
    "text": "%s"}`,
		from, from, to, subject, content,
	)

	url := "https://" + m.ApiURL + "/api/send"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(msg)))
	if err != nil {
		return fmt.Errorf("http.NewRequest(): error creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+m.Token)

	client := http.Client{Timeout: 15 * time.Second}

	res, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do(request): error sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("client.Do(request): status code not OK: %d", res.StatusCode)
	}

	/*
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("io.ReadAll: error reading response body: %w", err)
		}
		fmt.Println("body", string(body))
	*/

	return nil
}

func (m *MailTrap) SendEmailWithTemplate(
	ctx context.Context,
	from string,
	to string,
	templateUUID string,
	templateVars map[string]interface{},
) error {
	msg := map[string]interface{}{
		"from": map[string]string{
			"email": from,
			"name":  from,
		},
		"to": []map[string]string{
			{
				"email": to,
			},
		},
		"template_uuid":      templateUUID,
		"template_variables": templateVars,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}

	url := "https://" + m.ApiURL + "/api/send"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msgBytes))
	if err != nil {
		err = fmt.Errorf("http.NewRequest(): error creating request: %w", err)
		log.Println(err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+m.Token)

	client := http.Client{Timeout: 15 * time.Second}

	res, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("client.Do(request): error sending request: %w", err)
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("client.Do(request): status code not OK: %d", res.StatusCode)
		log.Println(err)
		return err
	}

	return nil
}
