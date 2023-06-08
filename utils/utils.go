package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net/mail"
	"net/smtp"
	"os"
	"regexp"
	"stark/database"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/palantir/stacktrace"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccessDetails struct {
	AccessUuid string
	UserID     string
}

type TokenDetail struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	AccessUuid     string `json:"access_uuid"`
	RefreshUuid    string `json:"refresh_uuid"`
	AccessExpires  int64  `json:"access_expires"`
	RefreshExpires int64  `json:"refresh_expires"`
}

type RefreshDetails struct {
	RefreshUuid string
	UserID      string
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Contains(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}

func GetJSONConfig() string {
	return `
	{
		"type": "` + os.Getenv("GOOGLE_CRED_TYPE") + `",
		"project_id": "` + os.Getenv("GOOGLE_CRED_PROJECT_ID") + `",
		"private_key_id": "` + os.Getenv("GOOGLE_CRED_PRIVATE_KEY_ID") + `",
		"private_key": "` + os.Getenv("GOOGLE_CRED_PRIVATE_KEY") + `",
		"client_email": "` + os.Getenv("GOOGLE_CRED_CLIENT_EMAIL") + `",
		"client_id": "` + os.Getenv("GOOGLE_CRED_CLIENT_ID") + `",
		"auth_uri": "` + os.Getenv("GOOGLE_CRED_AUTH_URI") + `",
		"token_uri": "` + os.Getenv("GOOGLE_CRED_TOKEN_URI") + `",
		"auth_provider_x509_cert_url": "` + os.Getenv("GOOGLE_CRED_AUTH_PROVIDER_X509_CERT_URL") + `",
		"client_x509_cert_url": "` + os.Getenv("GOOGLE_CRED_CLIENT_X509_CERT_URL") + `"
  	}`
}

func IsInList(list []string, s string) bool {
	for _, str := range list {
		if str == s {
			return true
		}
	}
	return false
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func SplitToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GetBearerKey(bearerToken string) (string, error) {
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) < 2 {
		return "", errors.New("invalid bearer token")
	}

	return strArr[1], nil
}

func CreateToken(user_id string) (*TokenDetail, error) {
	tokenDetail := &TokenDetail{}
	tokenDetail.AccessExpires = time.Now().Add(time.Hour * 24 * 30).Unix()
	tokenDetail.AccessUuid = uuid.NewV4().String()

	tokenDetail.RefreshExpires = time.Now().Add(time.Hour * 24 * 365).Unix()
	tokenDetail.RefreshUuid = tokenDetail.AccessUuid + "++" + user_id

	var err error
	// Creating access token
	accessClaims := jwt.MapClaims{}
	accessClaims["authorized"] = true
	accessClaims["access_uuid"] = tokenDetail.AccessUuid
	accessClaims["user_id"] = user_id
	accessClaims["exp"] = tokenDetail.AccessExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenDetail.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error when creating access token")
	}

	// Creating refresh token
	refreshClaims := jwt.MapClaims{}
	refreshClaims["refresh_uuid"] = tokenDetail.RefreshUuid
	refreshClaims["user_id"] = user_id
	refreshClaims["exp"] = tokenDetail.RefreshExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	tokenDetail.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error when creating refresh token")
	}

	return tokenDetail, nil
}

func ExtractAccessTokenMetadata(bearerToken string) (*AccessDetails, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		return &AccessDetails{
			AccessUuid: accessUuid,
			UserID:     fmt.Sprintf("%s", claims["user_id"]),
		}, nil
	}

	return nil, err
}

func FetchAccessAuth(authDetails *AccessDetails, redisDB *database.Redis) (string, error) {
	userID, err := redisDB.Get(authDetails.AccessUuid)
	if err != nil {
		return "", errors.New("Token expired")
	}

	if authDetails.UserID != userID {
		return "", errors.New("User not match")
	}

	return userID, nil
}

func ExtractRefreshTokenMetadata(bearerToken string) (*RefreshDetails, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, err
		}

		return &RefreshDetails{
			RefreshUuid: refreshUuid,
			UserID:      fmt.Sprintf("%s", claims["user_id"]),
		}, nil
	}

	return nil, err
}

func FetchRefreshAuth(authDetails *RefreshDetails, redisDB *database.Redis) (string, error) {
	userID, err := redisDB.Get(authDetails.RefreshUuid)
	if err != nil {
		return "", errors.New("token expired")
	}

	if authDetails.UserID != userID {
		return "", errors.New("user not match")
	}

	return userID, nil
}

func GetErrorMessage(validationError validator.FieldError) string {
	switch validationError.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "this field must be an email"
	case "lte":
		return "should be less than " + validationError.Param()
	case "gte":
		return "should be greater than " + validationError.Param()
	}

	return "unknown error"
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func SendMail(to []string, cc []string, subject, message string) error {
	from := mail.Address{Name: "Gimsak", Address: os.Getenv("SMTP_SENDER")}
	body := "From: " + from.String() + "\r\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		message

	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))
	smtpAddr := fmt.Sprintf("%s:%s", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))

	err := smtp.SendMail(smtpAddr, auth, os.Getenv("SMTP_EMAIL"), append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
