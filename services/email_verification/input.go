package email_verification

type Input struct {
	Email string `json:"email" binding:"required,email"`
}
