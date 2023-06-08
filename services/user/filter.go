package user

type Filter struct {
	Emails    []string `json:"emails"`
	Usernames []string `json:"usernames"`
}

func (f Filter) IsEmpty() bool {
	return len(f.Emails) == 0 && len(f.Usernames) == 0
}
