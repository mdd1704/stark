package user

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

func (s *Service) Create(name, email, username, contact, password string) (*User, error) {
	item := New(name, email, username, contact, password)
	totalByEmail, err := s.repo.FindTotalByFilter(Filter{Emails: []string{email}})
	if err != nil {
		return nil, err
	}

	totalByUsername, err := s.repo.FindTotalByFilter(Filter{Usernames: []string{username}})
	if err != nil {
		return nil, err
	}

	if totalByEmail > 0 || totalByUsername > 0 {
		return nil, failure.WithMessage(
			failure.CodeUserAlreadyExist,
			"username or email exists, duplicate username or email is not allowed",
		)
	}

	err = s.repo.Store(item)
	if err != nil {
		return nil, err
	}

	return s.FindByID(item.ID)
}

func (s *Service) Update(id uuid.UUID, name, email, username, contact, password string) (*User, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if stacktrace.RootCause(err) == sql.ErrNoRows {
			return nil, failure.WithMessage(
				failure.CodeUserNotFound,
				"user not found, id isn't in database",
			)
		}

		return nil, err
	}

	item.Update(name, email, username, contact, password)
	err = s.repo.Store(item)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *Service) UpdateProfile(id uuid.UUID, name, username, contact string) (*User, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if stacktrace.RootCause(err) == sql.ErrNoRows {
			return nil, failure.WithMessage(
				failure.CodeUserNotFound,
				"user not found, id isn't in database",
			)
		}

		return nil, err
	}

	item.UpdateProfile(name, username, contact)
	err = s.repo.StoreProfile(item)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *Service) FindByID(id uuid.UUID) (*User, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if stacktrace.RootCause(err) == sql.ErrNoRows {
			return nil, failure.WithMessage(
				failure.CodeUserNotFound,
				"user not found, id isn't in database",
			)
		}

		return nil, err
	}

	return item, nil
}

func (s *Service) FindAllByFilter(filter Filter) ([]*User, error) {
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
