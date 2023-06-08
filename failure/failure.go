package failure

import "fmt"

type Failure struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func New(code string) Failure {
	return Failure{Code: code}
}

func WithMessage(code string, desc string) Failure {
	return Failure{Code: code, Desc: desc}
}

func (f Failure) Error() string {
	return fmt.Sprintf("%s: %s", f.Code, f.Desc)
}

type Code string

const (
	CodeInternal           Code = "Internal"
	CodeClientNotFound          = "ClientNotFound"
	CodeUserAlreadyExist        = "UserAlreadyExist"
	CodeUserNotFound            = "UserNotFound"
	CodeUserDetailNotFound      = "UserDetailNotFound"
	CodeLoginFailed             = "LoginFailed"
	CodeIncorrectPassword       = "IncorrectPassword"
	CodeIncorrectToken          = "IncorrectToken"
	CodeTokenExpired            = "TokenExpired"
	CodeUserNotMatch            = "UserNotMatch"
	CodeIncorrectUserID         = "IncorrectUserID"
	CodeTokenAlreadyExist       = "TokenAlreadyExist"
)
