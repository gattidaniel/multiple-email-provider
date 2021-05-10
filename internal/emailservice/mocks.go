package emailservice

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/multiple-email-provider/cmd/server/requests"
	"github.com/stretchr/testify/mock"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) Translate(email requests.SendEmail) (SendEmailInput, error) {
	return m.Called(email).Get(0).(SendEmailInput), m.Called(email).Error(1)
}

func (m *MockEmailService) SendEmail(email SendEmailInput) error {
	return m.Called(email).Error(0)
}

// MockHttpClient is the mock client
type MockHttpClient struct {
}

// Do is the mock client's `Do` func
func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	// Maybe this is not the best approach, but it is the fastest
	switch req.Header.Get("X-Postmark-Server-Token") {
	case "requestError":
		return &http.Response{}, errors.New("error")
	case "badRequest":
		r := ioutil.NopCloser(bytes.NewReader([]byte(`{"error":"error"}`)))
		return &http.Response{StatusCode: http.StatusBadRequest, Body: r}, nil
	case "accepted":
		return &http.Response{StatusCode: http.StatusAccepted}, nil
	}

	return &http.Response{}, errors.New("not expected")
}
