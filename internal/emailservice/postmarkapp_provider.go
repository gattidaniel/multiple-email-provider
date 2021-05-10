package emailservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PostmarkProvider struct {
	ApiKey     string
	HTTPClient HTTPClient
}

func NewPostmarkProvider(apikey string, httpClient HTTPClient) PostmarkProvider {
	return PostmarkProvider{
		ApiKey:     apikey,
		HTTPClient: httpClient,
	}
}

type PostmarkInput struct {
	From          string `json:"From"`
	To            string `json:"To"`
	Subject       string `json:"Subject"`
	HtmlBody      string `json:"HtmlBody"`
	MessageStream string `json:"MessageStream"`
}

func (p PostmarkProvider) sendEmail(sendEmail SendEmailInput) error {
	url := "https://api.postmarkapp.com/email"

	sendGripInput := sendEmail.generatePostmarkInput()
	jsonValue, err := json.Marshal(sendGripInput)
	if err != nil {
		return fmt.Errorf("fail to marshal input %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("new request fail: %w", err)
	}
	req.Header.Add("X-Postmark-Server-Token", p.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.HTTPClient.Do(req)
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

	return fmt.Errorf("error with postmarkapp: Status code %d. Body %s", resp.StatusCode, string(body))
}

func (sendEmail SendEmailInput) generatePostmarkInput() PostmarkInput {
	return PostmarkInput{
		From:          sendEmail.From,
		To:            sendEmail.To,
		Subject:       sendEmail.Subject,
		HtmlBody:      sendEmail.HtmlBody,
		MessageStream: "outbound",
	}
}
