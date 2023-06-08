package user_detail

type Input struct {
	ID          string `json:"id" binding:"required"`
	DeviceToken string `json:"device_token"`
	DeviceOS    string `json:"device_os"`
	AvatarUrl   string `json:"avatar_url"`
	AvatarPath  string `json:"avatar_path"`
	Source      string `json:"source" binding:"required"`
	OAuthId     string `json:"oauth_id"`
	IDCardUrl   string `json:"id_card_url"`
	IDCardPath  string `json:"id_card_path"`
}
