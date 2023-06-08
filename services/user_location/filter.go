package user_location

type Filter struct {
	IDs []string `json:"ids"`
}

func (f Filter) IsEmpty() bool {
	return len(f.IDs) == 0
}
