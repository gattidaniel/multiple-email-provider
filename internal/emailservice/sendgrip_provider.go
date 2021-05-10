package emailservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SendGripProvider struct {
	ApiKey     string
	HTTPClient HTTPClient
}

func NewSendGripProvider(apikey string, httpClient HTTPClient) SendGripProvider {
	return SendGripProvider{
		ApiKey:     apikey,
		HTTPClient: httpClient,
	}
}

type SendGripInput struct {
	Personalization []Personalization `json:"personalizations"`
	Content         []Content         `json:"content"`
	From            Email             `json:"from"`
	ReplyTo         Email             `json:"reply_to"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Personalization struct {
	To      []Email `json:"to"`
	Subject string  `json:"subject"`
}
type Email struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (p SendGripProvider) sendEmail(sendEmail SendEmailInput) error {
	url := "https://api.sendgrid.com/v3/mail/send"

	sendGripInput := sendEmail.generateSendGripInput()
	jsonValue, err := json.Marshal(sendGripInput)
	if err != nil {
		return fmt.Errorf("fail to marshal input %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("new request fail: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+p.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request fail: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}

	return fmt.Errorf("error with sendgrip: Status code %d. Body %s", resp.StatusCode, string(body))
}

func (sendEmail SendEmailInput) generateSendGripInput() SendGripInput {
	return SendGripInput{
		Personalization: []Personalization{{
			To: []Email{{
				Email: sendEmail.To,
				Name:  sendEmail.ToName,
			}},
			Subject: sendEmail.Subject,
		}},
		Content: []Content{{
			Type:  "text/plain",
			Value: sendEmail.PlainBody,
		}},
		From: Email{
			Email: sendEmail.From,
			Name:  sendEmail.FromName,
		},
		ReplyTo: Email{
			Email: sendEmail.From,
			Name:  sendEmail.FromName,
		},
	}
}
