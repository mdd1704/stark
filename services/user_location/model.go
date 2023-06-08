package user_location

import (
	"time"

	"github.com/google/uuid"
)

type UserLocation struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ProvinceID string    `json:"province_id" db:"province_id"`
	RegencyID  string    `json:"regency_id" db:"regency_id"`
	DistrictID string    `json:"district_id" db:"district_id"`
	VillageID  string    `json:"village_id" db:"village_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func New(
	id uuid.UUID,
	province_id,
	regency_id,
	district_id,
	village_id string) *UserLocation {
	return &UserLocation{
		ID:         id,
		ProvinceID: province_id,
		RegencyID:  regency_id,
		DistrictID: district_id,
		VillageID:  village_id,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (u *UserLocation) Update(
	province_id,
	regency_id,
	district_id,
	village_id string) {
	u.ProvinceID = province_id
	u.RegencyID = regency_id
	u.DistrictID = district_id
	u.VillageID = village_id
	u.UpdatedAt = time.Now()
}

type Page struct {
	Items []*UserLocation `json:"items"`
	Total int             `json:"total"`
}
