package auth

import (
	"errors"
	"stark/database"
	"stark/failure"
	"stark/services/email_verification"
	"stark/services/user"
	"stark/utils"
	"strings"
	"time"

	"github.com/palantir/stacktrace"
)

const (
	emailVerificationSubject = "Email Verification"
	emailVerificationPreview = "Verifikasi email akun Gimsak kamu!"
)

type Service struct {
	redisDB                  *database.Redis
	userService              *user.Service
	emailVerificationService *email_verification.Service
}

func NewService(
	redisDB *database.Redis,
	userService *user.Service,
	emailVerificationService *email_verification.Service,
) *Service {
	return &Service{
		redisDB:                  redisDB,
		userService:              userService,
		emailVerificationService: emailVerificationService,
	}
}

func (s *Service) Login(email, username, password string) (*Login, error) {
	filter := user.Filter{}
	if email != "" {
		filter.Emails = []string{email}
	}

	if username != "" {
		filter.Usernames = []string{username}
	}

	if filter.IsEmpty() {
		return nil, failure.WithMessage(
			failure.CodeLoginFailed,
			"login failed, email or username cannot be empty",
		)
	}

	user, err := s.userService.FindAllByFilter(filter)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user[0].Password) {
		return nil, failure.WithMessage(
			failure.CodeIncorrectPassword,
			"incorrect password, try again",
		)
	}

	token, err := utils.CreateToken(user[0].ID.String())
	if err != nil {
		return nil, err
	}

	accessExpires := time.Unix(token.AccessExpires, 0)
	refreshExpires := time.Unix(token.RefreshExpires, 0)
	now := time.Now()

	err = s.redisDB.Set(token.AccessUuid, user[0].ID.String(), accessExpires.Sub(now))
	if err != nil {
		return nil, err
	}

	err = s.redisDB.Set(token.RefreshUuid, user[0].ID.String(), refreshExpires.Sub(now))
	if err != nil {
		return nil, err
	}

	login := &Login{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return login, nil
}

func (s *Service) Logout(accessUuid string, userID string) (int64, error) {
	deleted, err := s.redisDB.Delete(accessUuid)
	if err != nil {
		return 0, err
	}

	deleted, err = s.redisDB.Delete(accessUuid + "++" + userID)
	if err != nil {
		return 0, err
	}

	return deleted, nil
}

func (s *Service) RefreshToken(refreshToken string) (*Login, error) {
	metadata, err := utils.ExtractRefreshTokenMetadata(refreshToken)
	if err != nil {
		return nil, failure.WithMessage(
			failure.CodeIncorrectToken,
			"incorrect token, try again",
		)
	}

	userID, err := utils.FetchRefreshAuth(metadata, s.redisDB)
	if err != nil {
		if stacktrace.RootCause(err).Error() == "token expired" {
			return nil, failure.WithMessage(
				failure.CodeTokenExpired,
				"token expired, need to login again",
			)
		}

		if stacktrace.RootCause(err).Error() == "user not match" {
			return nil, failure.WithMessage(
				failure.CodeUserNotMatch,
				"user not match, invalid token",
			)
		}

		return nil, err
	}

	splitRefreshUuid := strings.Split(metadata.RefreshUuid, "++")
	if len(splitRefreshUuid) != 2 {
		return nil, errors.New("invalid refresh uuid")
	}

	_, err = s.redisDB.Delete(metadata.RefreshUuid)
	if err != nil {
		return nil, err
	}

	_, err = s.redisDB.Delete(splitRefreshUuid[0])
	if err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(userID)
	if err != nil {
		return nil, err
	}

	accessExpires := time.Unix(token.AccessExpires, 0)
	refreshExpires := time.Unix(token.RefreshExpires, 0)
	now := time.Now()

	err = s.redisDB.Set(token.AccessUuid, userID, accessExpires.Sub(now))
	if err != nil {
		return nil, err
	}

	err = s.redisDB.Set(token.RefreshUuid, userID, refreshExpires.Sub(now))
	if err != nil {
		return nil, err
	}

	login := &Login{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return login, nil
}

func (s *Service) Register(name, email, username, contact, password string) error {
	_, err := s.userService.Create(name, email, username, contact, password)
	if err != nil {
		return err
	}

	emailVerification, err := s.emailVerificationService.Create(email)
	if err != nil {
		return err
	}

	to := []string{email}
	content := verificationEmailContent(email, emailVerification.Token)
	message := utils.EmailLayout(emailVerificationPreview, content)
	err = utils.SendMail(to, nil, emailVerificationSubject, message)
	if err != nil {
		return err
	}

	return nil
}
