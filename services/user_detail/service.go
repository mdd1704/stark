package user_detail

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
	device_token,
	device_os,
	avatar_url,
	avatar_path,
	source,
	oauth_id,
	id_card_url,
	id_card_path string) (*UserDetail, error) {
	item := New(
		id,
		device_token,
		device_os,
		avatar_url,
		avatar_path,
		source,
		oauth_id,
		id_card_url,
		id_card_path,
	)

	var err error
	totalOauthID := 0
	if oauth_id != "" {
		totalOauthID, err = s.repo.FindTotalByFilter(Filter{OAuthIDs: []string{oauth_id}})
		if err != nil {
			return nil, err
		}
	}

	totalID, err := s.repo.FindTotalByFilter(Filter{IDs: []string{id.String()}})
	if err != nil {
		return nil, err
	}

	if totalOauthID > 0 || totalID > 0 {
		return nil, failure.WithMessage(
			failure.CodeUserDetailNotFound,
			"id or oauth_id exists, duplicate id or oauth_id is not allowed",
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
	device_token,
	device_os,
	avatar_url,
	avatar_path,
	source,
	oauth_id,
	id_card_url,
	id_card_path string) (*UserDetail, error) {
	item := &UserDetail{
		ID:          id,
		DeviceToken: device_token,
		DeviceOS:    device_os,
		AvatarUrl:   avatar_url,
		AvatarPath:  avatar_path,
		Source:      source,
		OAuthId:     oauth_id,
		IDCardUrl:   id_card_url,
		IDCardPath:  id_card_path,
	}

	item.Update(
		device_token,
		device_os,
		avatar_url,
		avatar_path,
		source,
		oauth_id,
		id_card_url,
		id_card_path,
	)

	err := s.repo.Store(item)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *Service) FindByID(id uuid.UUID) (*UserDetail, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if stacktrace.RootCause(err) == sql.ErrNoRows {
			return nil, failure.WithMessage(
				failure.CodeUserNotFound,
				"user detail not found, id isn't in database",
			)
		}

		return nil, err
	}

	return item, nil
}

func (s *Service) FindAllByFilter(filter Filter) ([]*UserDetail, error) {
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
