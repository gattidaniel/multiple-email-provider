package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/multiple-email-provider/internal/emailservice"

	"github.com/multiple-email-provider/cmd/server/httphandlers"
)

func main() {
	postmarkKey, sendGripKey := getKeyFromEnvironment()
	httpClient := http.Client{}

	postmarkProvider := emailservice.NewPostmarkProvider(postmarkKey, &httpClient)
	sendGripProvider := emailservice.NewPostmarkProvider(sendGripKey, &httpClient)

	emailService, err := emailservice.NewEmailService(postmarkProvider, sendGripProvider)
	if err != nil {
		log.Fatal(err)
	}

	emailHandler := httphandlers.NewEmailHandler(emailService)
	http.HandleFunc("/v1/email", emailHandler.Send)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getKeyFromEnvironment() (string, string) {
	postmarkKey := os.Getenv("postmark-key")
	if postmarkKey == "" {
		log.Fatal(errors.New("set postmark-key with appropiate key"))
	}
	sendGripKey := os.Getenv("sendgrip-key")
	if sendGripKey == "" {
		log.Fatal(errors.New("set postmark-key with appropiate key"))
	}
	return postmarkKey, sendGripKey
}
