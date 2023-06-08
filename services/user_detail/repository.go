package user_detail

import "github.com/google/uuid"

type Repository interface {
	Store(item *UserDetail) error
	FindByID(id uuid.UUID) (*UserDetail, error)
	FindByFilter(filter Filter) ([]*UserDetail, error)
	FindPage(offset, limit int) ([]*UserDetail, error)
	FindTotalByFilter(filter Filter) (int, error)
}
