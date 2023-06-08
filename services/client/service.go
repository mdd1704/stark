package client

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

func (s *Service) Create(name string) (*Client, error) {
	item := New(name)
	for {
		total, err := s.repo.FindTotalByFilter(Filter{BearerKeys: []string{item.BearerKey}})
		if err != nil {
			return nil, err
		}

		if total != 0 {
			item = New(name)
			continue
		}

		break
	}

	err := s.repo.Store(item)
	if err != nil {
		return nil, err
	}

	return s.FindByID(item.ID)
}

func (s *Service) Update(id uuid.UUID, name string) (*Client, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if stacktrace.RootCause(err) == sql.ErrNoRows {
			return nil, failure.WithMessage(
				failure.CodeClientNotFound,
				"client not found, id isn't in database",
			)
		}

		return nil, err
	}

	item.Update(name)
	err = s.repo.Store(item)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *Service) FindByID(id uuid.UUID) (*Client, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if stacktrace.RootCause(err) == sql.ErrNoRows {
			return nil, failure.WithMessage(
				failure.CodeClientNotFound,
				"client not found, id isn't in database",
			)
		}

		return nil, err
	}

	return item, nil
}

func (s *Service) FindAllByFilter(filter Filter) ([]*Client, error) {
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
