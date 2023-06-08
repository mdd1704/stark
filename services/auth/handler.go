package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/palantir/stacktrace"

	"stark/failure"
	"stark/respond"
	"stark/utils"
	"stark/utils/activity"
	"stark/utils/log"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleLogin(c *gin.Context) {
	ctx := activity.NewContext("auth_login")
	trx, _ := activity.GetTransactionID(ctx)
	var input InputLogin

	if err := c.ShouldBindJSON(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			out := make([]utils.ErrorMessage, 0)
			for _, validationError := range validationErrors {
				out = append(out, utils.ErrorMessage{
					Field:   utils.ToSnakeCase(validationError.Field()),
					Message: utils.GetErrorMessage(validationError),
				})
			}

			respond.Invalid(c, trx, http.StatusBadRequest, out)
		}
		return
	}

	token, err := h.service.Login(input.Email, input.Username, input.Password)
	if err != nil {
		if f, ok := stacktrace.RootCause(err).(failure.Failure); ok {
			switch f.Code {
			case failure.CodeUserNotFound:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			case failure.CodeIncorrectPassword:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			}
		}

		log.WithContext(ctx).Error(stacktrace.Propagate(err, "auth login error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, "unknown error")
		return
	}

	respond.Success(c, trx, http.StatusCreated, token)
}

func (h *Handler) HandleRefreshToken(c *gin.Context) {
	ctx := activity.NewContext("auth_refresh_token")
	trx, _ := activity.GetTransactionID(ctx)
	var input InputRefreshToken

	if err := c.ShouldBindJSON(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			out := make([]utils.ErrorMessage, 0)
			for _, validationError := range validationErrors {
				out = append(out, utils.ErrorMessage{
					Field:   utils.ToSnakeCase(validationError.Field()),
					Message: utils.GetErrorMessage(validationError),
				})
			}

			respond.Invalid(c, trx, http.StatusBadRequest, out)
		}
		return
	}

	token, err := h.service.RefreshToken(input.RefreshToken)
	if err != nil {
		if f, ok := stacktrace.RootCause(err).(failure.Failure); ok {
			switch f.Code {
			case failure.CodeIncorrectToken:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			case failure.CodeTokenExpired:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			case failure.CodeUserNotMatch:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			}
		}

		log.WithContext(ctx).Error(stacktrace.Propagate(err, "auth refresh token error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, "unknown error")
		return
	}

	respond.Success(c, trx, http.StatusCreated, token)
}

func (h *Handler) HandleLogout(c *gin.Context) {
	ctx := activity.NewContext("auth_logout")
	ctx = activity.WithUserID(ctx, c.Value("user_id").(string))
	trx, _ := activity.GetTransactionID(ctx)
	userID, _ := activity.GetUserID(ctx)
	accessUuid := c.Value("access_uuid").(string)

	token, err := h.service.Logout(accessUuid, userID)
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "auth logout error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, "unknown error")
		return
	}

	respond.Success(c, trx, http.StatusCreated, token)
}

func (h *Handler) HandleRegister(c *gin.Context) {
	ctx := activity.NewContext("auth_register")
	trx, _ := activity.GetTransactionID(ctx)
	var input InputRegister

	if err := c.ShouldBindJSON(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			out := make([]utils.ErrorMessage, 0)
			for _, validationError := range validationErrors {
				out = append(out, utils.ErrorMessage{
					Field:   utils.ToSnakeCase(validationError.Field()),
					Message: utils.GetErrorMessage(validationError),
				})
			}

			respond.Invalid(c, trx, http.StatusBadRequest, out)
		}
		return
	}

	err := h.service.Register(input.Name, input.Email, input.Username, input.Contact, input.Password)
	if err != nil {
		if f, ok := stacktrace.RootCause(err).(failure.Failure); ok {
			switch f.Code {
			case failure.CodeUserAlreadyExist:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			}
		}

		log.WithContext(ctx).Error(stacktrace.Propagate(err, "auth register error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, "unknown error")
		return
	}

	respond.Success(c, trx, http.StatusCreated, nil)
}
