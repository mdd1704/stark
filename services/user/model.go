package user

import (
	"time"

	"github.com/google/uuid"

	"stark/utils"
)

type User struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Email           string    `json:"email" db:"email"`
	Username        string    `json:"username" db:"username"`
	Contact         string    `json:"contact" db:"contact"`
	Password        string    `json:"password" db:"password"`
	EmailVerifiedAt time.Time `json:"email_verified_at" db:"email_verified_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func New(name, email, username, contact, password string) *User {
	id := uuid.New()
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		hashPassword = "-"
	}

	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		Username:  username,
		Contact:   contact,
		Password:  hashPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) Update(name, email, username, contact, password string) {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		hashPassword = "-"
	}

	u.Name = name
	u.Email = email
	u.Username = username
	u.Contact = contact
	u.Password = hashPassword
	u.UpdatedAt = time.Now()
}

func (u *User) UpdateProfile(name, username, contact string) {
	u.Name = name
	u.Username = username
	u.Contact = contact
	u.UpdatedAt = time.Now()
}

type Page struct {
	Items []*User `json:"items"`
	Total int     `json:"total"`
}
