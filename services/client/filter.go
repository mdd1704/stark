package client

type Filter struct {
	Names      []string `json:"names"`
	BearerKeys []string `json:"bearer_keys"`
}

func (f Filter) IsEmpty() bool {
	return len(f.Names) == 0 && len(f.BearerKeys) == 0
}
