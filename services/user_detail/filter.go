package user_detail

type Filter struct {
	IDs          []string `json:"ids"`
	OAuthIDs     []string `json:"oauth_ids"`
	DeviceTokens []string `json:"device_tokens"`
	Sources      []string `json:"sources"`
}

func (f Filter) IsEmpty() bool {
	return len(f.IDs) == 0 && len(f.OAuthIDs) == 0 && len(f.DeviceTokens) == 0 && len(f.Sources) == 0
}
