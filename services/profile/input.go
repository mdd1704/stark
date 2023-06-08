package profile

type InputLogin struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type InputRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type InputUpdateProfile struct {
	Name        string `json:"name"`
	Username    string `json:"username" binding:"required,max=25"`
	Contact     string `json:"contact"`
	DeviceToken string `json:"device_token"`
	DeviceOS    string `json:"device_os"`
	AvatarUrl   string `json:"avatar_url"`
	AvatarPath  string `json:"avatar_path"`
	Source      string `json:"source"`
	OAuthId     string `json:"oauth_id"`
	IDCardUrl   string `json:"id_card_url"`
	IDCardPath  string `json:"id_card_path"`
	ProvinceID  string `json:"province_id"`
	RegencyID   string `json:"regency_id"`
	DistrictID  string `json:"district_id"`
	VillageID   string `json:"village_id"`
}

type InputChangePassword struct {
	OldPassword             string `json:"old_password" binding:"required"`
	NewPassword             string `json:"new_password" binding:"required"`
	NewPasswordConfirmation string `json:"new_password_confirmation" binding:"required"`
}
