package user

import "github.com/google/uuid"

type Repository interface {
	Store(data *User) error
	StoreProfile(data *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByFilter(filter Filter) ([]*User, error)
	FindPage(offset, limit int) ([]*User, error)
	FindTotalByFilter(filter Filter) (int, error)
}
