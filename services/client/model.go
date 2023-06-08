package client

import (
	"time"

	"github.com/google/uuid"

	"stark/utils"
)

type Client struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	BearerKey string    `json:"bearer_key" db:"bearer_key"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func New(name string) *Client {
	id := uuid.New()
	bearer_key := utils.GenerateSecureToken(25)

	return &Client{
		ID:        id,
		Name:      name,
		BearerKey: bearer_key,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *Client) Update(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

type Page struct {
	Items []*Client `json:"items"`
	Total int       `json:"total"`
}
