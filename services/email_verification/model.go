package email_verification

import (
	"time"

	"stark/utils"
)

type EmailVerification struct {
	Email     string    `json:"email" db:"email"`
	Token     string    `json:"token" db:"token"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func New(email string) *EmailVerification {
	token := utils.GenerateSecureToken(25)

	return &EmailVerification{
		Email:     email,
		Token:     token,
		CreatedAt: time.Now(),
	}
}
