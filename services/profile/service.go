package profile

import (
	"github.com/google/uuid"

	"stark/failure"
	"stark/services/user"
	"stark/services/user_detail"
	"stark/services/user_location"
	"stark/utils"
)

type Service struct {
	userService         *user.Service
	userDetailService   *user_detail.Service
	userLocationService *user_location.Service
}

func NewService(
	userService *user.Service,
	userDetailService *user_detail.Service,
	userLocationService *user_location.Service,
) *Service {
	return &Service{
		userService:         userService,
		userDetailService:   userDetailService,
		userLocationService: userLocationService,
	}
}

func (s *Service) UpdateProfile(
	id uuid.UUID,
	name,
	username,
	contact,
	device_token,
	device_os,
	avatar_url,
	avatar_path,
	source,
	oauth_id,
	id_card_url,
	id_card_path,
	province_id,
	regency_id,
	district_id,
	village_id string) error {
	_, err := s.userService.UpdateProfile(
		id,
		name,
		username,
		contact,
	)

	if err != nil {
		return err
	}

	_, err = s.userDetailService.Update(
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

	if err != nil {
		return err
	}

	_, err = s.userLocationService.Update(
		id,
		province_id,
		regency_id,
		district_id,
		village_id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ChangePassword(id uuid.UUID, old_password, new_password, new_password_confirmation string) error {
	user, err := s.userService.FindByID(id)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(old_password, user.Password) {
		return failure.WithMessage(
			failure.CodeIncorrectPassword,
			"incorrect password, try again",
		)
	}

	_, err = s.userService.Update(user.ID, user.Name, user.Email, user.Username, user.Contact, new_password)
	if err != nil {
		return err
	}

	return nil
}
