package emailservice

import "net/http"

type emailProvider interface {
	sendEmail(SendEmailInput) error
}
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
