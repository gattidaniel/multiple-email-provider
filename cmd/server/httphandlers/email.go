package httphandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/multiple-email-provider/cmd/server/requests"
	"github.com/multiple-email-provider/internal/emailservice"
	"github.com/multiple-email-provider/internal/validatorhelper"
)

type EmailHandler struct {
	EmailService EmailService
}

type EmailService interface {
	Translate(email requests.SendEmail) (emailservice.SendEmailInput, error)
	SendEmail(email emailservice.SendEmailInput) error
}

var (
	validate = validator.New()
)

func NewEmailHandler(emailService emailservice.EmailService) EmailHandler {
	return EmailHandler{
		EmailService: emailService,
	}
}

func (h EmailHandler) Send(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		http.Error(w, "request body is empty", http.StatusBadRequest)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	var sendEmailRequest requests.SendEmail
	if err := json.Unmarshal(bodyBytes, &sendEmailRequest); err != nil {
		http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = validate.Struct(sendEmailRequest)
	if err != nil {
		errorMessage := validatorhelper.ProcessErrors(err.(validator.ValidationErrors))
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	sendEmail, err := h.EmailService.Translate(sendEmailRequest)
	if err != nil {
		err = fmt.Errorf("fail to translate: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.EmailService.SendEmail(sendEmail)
	if err != nil {
		err = fmt.Errorf("fail to send email: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
