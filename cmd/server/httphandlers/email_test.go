package httphandlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/multiple-email-provider/cmd/server/requests"
	"github.com/multiple-email-provider/internal/emailservice"
	"github.com/stretchr/testify/assert"
)

var (
	requestTranslateMockFail = requests.SendEmail{
		To:       "to@mail.com",
		ToName:   "to_name",
		From:     "from@mail.com",
		FromName: "from_name",
		Subject:  "fail",
		Body:     "fail",
	}
	requestTranslateMockOkSendFail = requests.SendEmail{
		To:       "to@mail.com",
		ToName:   "to_name",
		From:     "from@mail.com",
		FromName: "from_name",
		Subject:  "translate ok, send fail",
		Body:     "translate ok, send fail",
	}
	requestTranslateMockOkSendOk = requests.SendEmail{
		To:       "to@mail.com",
		ToName:   "to_name",
		From:     "from@mail.com",
		FromName: "from_name",
		Subject:  "translate ok, send ok",
		Body:     "translate ok, send ok",
	}
	requestSendEmailMockFail = emailservice.SendEmailInput{
		To:        "to@mail.com",
		ToName:    "to_name",
		From:      "from@mail.com",
		FromName:  "from_name",
		Subject:   "send_fail",
		PlainBody: "send_fail",
		HtmlBody:  "send_fail",
	}
	requestSendEmailMockOk = emailservice.SendEmailInput{
		To:        "to@mail.com",
		ToName:    "to_name",
		From:      "from@mail.com",
		FromName:  "from_name",
		Subject:   "send_ok",
		PlainBody: "send_ok",
		HtmlBody:  "send_ok",
	}
)

func TestEmailHandler_Send(t *testing.T) {
	requestTranslateFailMarshall, err := json.Marshal(requestTranslateMockFail)
	assert.NoError(t, err)
	requestTranslateMockOkSendFailMarshall, err := json.Marshal(requestTranslateMockOkSendFail)
	assert.NoError(t, err)
	requestTranslateMockOkSendOkMarshall, err := json.Marshal(requestTranslateMockOkSendOk)
	assert.NoError(t, err)
	mockEmailService := new(emailservice.MockEmailService)
	mockEmailService.On("Translate", requestTranslateMockFail).Return(emailservice.SendEmailInput{}, errors.New("error"))
	mockEmailService.On("Translate", requestTranslateMockOkSendFail).Return(requestSendEmailMockFail, nil)
	mockEmailService.On("Translate", requestTranslateMockOkSendOk).Return(requestSendEmailMockOk, nil)
	mockEmailService.On("SendEmail", requestSendEmailMockFail).Return(errors.New("error"))
	mockEmailService.On("SendEmail", requestSendEmailMockOk).Return(nil)

	type fields struct {
		EmailService EmailService
	}
	type args struct {
		method string
		body   string
	}
	type expected struct {
		responseCode int
		body         string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expected
	}{
		{
			name:   "wrong_method",
			fields: fields{},
			args: args{
				method: http.MethodGet,
				body:   "",
			},
			expected: expected{
				responseCode: http.StatusMethodNotAllowed,
				body:         "method not allowed",
			},
		},
		{
			name:   "body_empty",
			fields: fields{},
			args: args{
				method: http.MethodPost,
				body:   "",
			},
			expected: expected{
				responseCode: http.StatusBadRequest,
				body:         "error reading request body: unexpected end of JSON input",
			},
		},
		{
			name:   "body_json_invalid",
			fields: fields{},
			args: args{
				method: http.MethodPost,
				body:   "hi",
			},
			expected: expected{
				responseCode: http.StatusBadRequest,
				body:         "error reading request body: invalid character 'h' looking for beginning of value",
			},
		},
		{
			name: "validation_fail",
			fields: fields{
				EmailService: mockEmailService,
			},
			args: args{
				method: http.MethodPost,
				body:   "{}",
			},
			expected: expected{
				responseCode: http.StatusBadRequest,
				body:         "Error in field 'To'. Is required\nError in field 'ToName'. Is required\nError in field 'From'. Is required\nError in field 'FromName'. Is required\nError in field 'Subject'. Is required\nError in field 'Body'. Is required",
			},
		},
		{
			name: "translate_fail",
			fields: fields{
				EmailService: mockEmailService,
			},
			args: args{
				method: http.MethodPost,
				body:   string(requestTranslateFailMarshall),
			},
			expected: expected{
				responseCode: http.StatusInternalServerError,
				body:         "fail to translate: error",
			},
		},
		{
			name: "send_fail",
			fields: fields{
				EmailService: mockEmailService,
			},
			args: args{
				method: http.MethodPost,
				body:   string(requestTranslateMockOkSendFailMarshall),
			},
			expected: expected{
				responseCode: http.StatusInternalServerError,
				body:         "fail to send email: error",
			},
		},
		{
			name: "send_ok",
			fields: fields{
				EmailService: mockEmailService,
			},
			args: args{
				method: http.MethodPost,
				body:   string(requestTranslateMockOkSendOkMarshall),
			},
			expected: expected{
				responseCode: http.StatusOK,
				body:         "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := EmailHandler{
				EmailService: tt.fields.EmailService,
			}

			req, err := http.NewRequest(tt.args.method, "/email", strings.NewReader(tt.args.body))

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Send)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.expected.responseCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expected.responseCode)
			}

			// Check the response body is what we expect.
			bodyBytes, err := ioutil.ReadAll(rr.Body)
			assert.NoError(t, err)
			if body := string(bodyBytes); !strings.Contains(body, tt.expected.body) {
				t.Errorf("handler returned wrong body: got %v which not contain %v", body, tt.expected.body)
			}
		})
	}
}
