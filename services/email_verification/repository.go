package email_verification

type Repository interface {
	Store(data *EmailVerification) error
	FindByToken(token string) (*EmailVerification, error)
	FindTotalByFilter(filter Filter) (int, error)
}
