package user_location

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/palantir/stacktrace"

	"stark/failure"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	id uuid.UUID,
	province_id,
	regency_id,
	district_id,
	village_id string) (*UserLocation, error) {
	item := New(
		id,
		province_id,
		regency_id,
		district_id,
		village_id,
	)

	totalID, err := s.repo.FindTotalByFilter(Filter{IDs: []string{id.String()}})
	if err != nil {
		return nil, err
	}

	if totalID > 0 {
		return nil, failure.WithMessage(
			failure.CodeUserDetailNotFound,
			"id exists, duplicate id is not allowed",
		)
	}

	err = s.repo.Store(item)
	if err != nil {
		return nil, err
	}

	return s.FindByID(item.ID)
}

func (s *Service) Update(
	id uuid.UUID,
	province_id,
	regency_id,
	district_id,
	village_id string) (*UserLocation, error) {
	item := &UserLocation{
		ID:         id,
		ProvinceID: province_id,
		RegencyID:  regency_id,
		DistrictID: district_id,
		VillageID:  village_id,
	}

	item.Update(
		province_id,
		regency_id,
		district_id,
		village_id,
	)

	err := s.repo.Store(item)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *Service) FindByID(id uuid.UUID) (*UserLocation, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if stacktrace.RootCause(err) == sql.ErrNoRows {
			return nil, failure.WithMessage(
				failure.CodeUserNotFound,
				"user location not found, id isn't in database",
			)
		}

		return nil, err
	}

	return item, nil
}

func (s *Service) FindAllByFilter(filter Filter) ([]*UserLocation, error) {
	return s.repo.FindByFilter(filter)
}

func (s *Service) FindPage(page, limit int) (Page, error) {
	offset := (page - 1) * limit
	items, err := s.repo.FindPage(offset, limit)
	if err != nil {
		return Page{}, err
	}

	total, err := s.repo.FindTotalByFilter(Filter{})
	if err != nil {
		return Page{}, err
	}

	return Page{
		Items: items,
		Total: total,
	}, nil
}
