package user_detail

import (
	"time"

	"github.com/google/uuid"
)

type UserDetail struct {
	ID          uuid.UUID `json:"id" db:"id"`
	DeviceToken string    `json:"device_token" db:"device_token"`
	DeviceOS    string    `json:"device_os" db:"device_os"`
	AvatarUrl   string    `json:"avatar_url" db:"avatar_url"`
	AvatarPath  string    `json:"avatar_path" db:"avatar_path"`
	Source      string    `json:"source" db:"source"`
	OAuthId     string    `json:"oauth_id" db:"oauth_id"`
	IDCardUrl   string    `json:"id_card_url" db:"id_card_url"`
	IDCardPath  string    `json:"id_card_path" db:"id_card_path"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func New(
	id uuid.UUID,
	device_token,
	device_os,
	avatar_url,
	avatar_path,
	source,
	oauth_id,
	id_card_url,
	id_card_path string) *UserDetail {
	return &UserDetail{
		ID:          id,
		DeviceToken: device_token,
		DeviceOS:    device_os,
		AvatarUrl:   avatar_url,
		AvatarPath:  avatar_path,
		Source:      source,
		OAuthId:     oauth_id,
		IDCardUrl:   id_card_url,
		IDCardPath:  id_card_path,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (u *UserDetail) Update(
	device_token,
	device_os,
	avatar_url,
	avatar_path,
	source,
	oauth_id,
	id_card_url,
	id_card_path string) {
	u.DeviceToken = device_token
	u.DeviceOS = device_os
	u.AvatarUrl = avatar_url
	u.AvatarPath = avatar_path
	u.Source = source
	u.OAuthId = oauth_id
	u.IDCardUrl = id_card_url
	u.IDCardPath = id_card_path
	u.UpdatedAt = time.Now()
}

type Page struct {
	Items []*UserDetail `json:"items"`
	Total int           `json:"total"`
}
