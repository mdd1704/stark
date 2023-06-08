package user_location

import "github.com/google/uuid"

type Repository interface {
	Store(item *UserLocation) error
	FindByID(id uuid.UUID) (*UserLocation, error)
	FindByFilter(filter Filter) ([]*UserLocation, error)
	FindPage(offset, limit int) ([]*UserLocation, error)
	FindTotalByFilter(filter Filter) (int, error)
}
