package email_verification

import (
	"database/sql"

	"github.com/palantir/stacktrace"

	"stark/failure"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(email string) (*EmailVerification, error) {
	item := New(email)
	for {
		total, err := s.repo.FindTotalByFilter(Filter{Tokens: []string{item.Token}})
		if err != nil {
			return nil, err
		}

		if total != 0 {
			item = New(email)
			continue
		}

		break
	}

	err := s.repo.Store(item)
	if err != nil {
		return nil, err
	}

	return s.FindByToken(item.Token)
}

func (s *Service) FindByToken(token string) (*EmailVerification, error) {
	item, err := s.repo.FindByToken(token)
	if err != nil {
		if stacktrace.RootCause(err) == sql.ErrNoRows {
			return nil, failure.WithMessage(
				failure.CodeUserNotFound,
				"email verification not found, token isn't in database",
			)
		}

		return nil, err
	}

	return item, nil
}
