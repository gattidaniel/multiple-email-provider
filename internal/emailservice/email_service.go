package emailservice

import (
	"errors"
	"fmt"
	"github.com/multiple-email-provider/cmd/server/requests"
	"log"

	"jaytaylor.com/html2text"
)

type EmailService struct {
	providers []emailProvider
}

type SendEmailInput struct {
	To        string
	ToName    string
	From      string
	FromName  string
	Subject   string
	HtmlBody  string
	PlainBody string
}

func NewEmailService(providers ...emailProvider) (EmailService, error) {
	if len(providers) == 0 {
		return EmailService{}, errors.New("at least one provider is required")
	}
	return EmailService{providers: providers}, nil
}

func (s EmailService) Translate(email requests.SendEmail) (SendEmailInput, error) {
	plainBody, err := html2text.FromString(email.Body, html2text.Options{PrettyTables: true})
	if err != nil {
		return SendEmailInput{}, fmt.Errorf("error getting text from body: %w", err)
	}

	return SendEmailInput{
		To:        email.To,
		ToName:    email.ToName,
		From:      email.From,
		FromName:  email.FromName,
		Subject:   email.Subject,
		HtmlBody:  email.Body,
		PlainBody: plainBody,
	}, nil
}

func (s EmailService) SendEmail(email SendEmailInput) error {
	for _, provider := range s.providers {
		err := provider.sendEmail(email)
		log.Println(err.Error())
		// TODO: Improve error management
		// If we have a connection issue on our side, maybe mail was sent, but we didn't receive the answer of the
		// provider, so for the moment we can send a mail twice from different providers.
		if err == nil {
			return nil
		}
	}
	return errors.New("mail wasn't sent. Any provider could send it")
}
