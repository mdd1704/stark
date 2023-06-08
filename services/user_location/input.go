package user_location

type Input struct {
	ID         string `json:"id" binding:"required"`
	ProvinceID string `json:"province_id"`
	RegencyID  string `json:"regency_id"`
	DistrictID string `json:"district_id"`
	VillageID  string `json:"village_id"`
}
