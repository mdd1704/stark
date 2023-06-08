package email_verification

type Filter struct {
	Tokens []string `json:"tokens"`
}

func (f Filter) IsEmpty() bool {
	return len(f.Tokens) == 0
}
