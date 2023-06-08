package client

import "github.com/google/uuid"

type Repository interface {
	Store(item *Client) error
	FindByID(id uuid.UUID) (*Client, error)
	FindByFilter(filter Filter) ([]*Client, error)
	FindPage(offset, limit int) ([]*Client, error)
	FindTotalByFilter(filter Filter) (int, error)
}
