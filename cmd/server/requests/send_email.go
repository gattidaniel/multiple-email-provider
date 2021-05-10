package requests

type SendEmail struct {
	To       string `json:"to" validate:"required,email"`
	ToName   string `json:"to_name" validate:"required"`
	From     string `json:"from" validate:"required,email"`
	FromName string `json:"from_name" validate:"required"`
	Subject  string `json:"subject" validate:"required"`
	Body     string `json:"body" validate:"required"`
}
