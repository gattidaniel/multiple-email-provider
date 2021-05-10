package emailservice

import (
	"reflect"
	"testing"
)

func TestNewEmailService(t *testing.T) {
	type args struct {
		providers []emailProvider
	}
	tests := []struct {
		name    string
		args    args
		want    EmailService
		wantErr bool
	}{
		{
			name:    "providers_empty",
			args:    args{},
			want:    EmailService{},
			wantErr: true,
		},
		{
			name: "providers_ok",
			args: args{
				providers: []emailProvider{PostmarkProvider{ApiKey: ""}},
			},
			want:    EmailService{providers: []emailProvider{PostmarkProvider{ApiKey: ""}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmailService(tt.args.providers...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmailService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmailService() got = %v, want %v", got, tt.want)
			}
		})
	}
}
