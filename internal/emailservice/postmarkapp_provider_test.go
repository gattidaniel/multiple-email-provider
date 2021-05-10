package emailservice

import (
	"strings"
	"testing"
)

var (
	sendEmailInput = SendEmailInput{
		To:        "To",
		ToName:    "ToName",
		From:      "From",
		FromName:  "FromName",
		Subject:   "Error",
		HtmlBody:  "Error",
		PlainBody: "Error",
	}
)

func TestPostmarkProvider_sendEmail(t *testing.T) {
	mockHttpClient := new(MockHttpClient)

	type args struct {
		sendEmail SendEmailInput
		apiKey    string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr string
	}{
		{
			name: "RequestError",
			args: args{
				sendEmail: sendEmailInput,
				apiKey:    "requestError",
			},
			wantErr:     true,
			expectedErr: "do request fail: error",
		},
		{
			name: "BadRequest",
			args: args{
				sendEmail: sendEmailInput,
				apiKey:    "badRequest",
			},
			wantErr:     true,
			expectedErr: "error with postmarkapp: Status code 400. Body {\"error\":\"error\"}",
		},
		{
			name: "OK",
			args: args{
				sendEmail: sendEmailInput,
				apiKey:    "accepted",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPostmarkProvider(tt.args.apiKey, mockHttpClient)

			err := p.sendEmail(tt.args.sendEmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("sendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
