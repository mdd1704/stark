package profile

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (h *Handler) HandleUpdateProfile(c *gin.Context) {
	ctx := activity.NewContext("update_profile")
	ctx = activity.WithUserID(ctx, c.Value("user_id").(string))
	trx, _ := activity.GetTransactionID(ctx)
	userID, _ := activity.GetUserID(ctx)
	var input InputUpdateProfile

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

	id, err := uuid.Parse(userID)
	if err != nil {
		respond.Error(c, trx, http.StatusBadRequest, respond.ErrBadRequest, "invalid user id")
		return
	}

	err = h.service.UpdateProfile(
		id,
		input.Name,
		input.Username,
		input.Contact,
		input.DeviceToken,
		input.DeviceOS,
		input.AvatarUrl,
		input.AvatarPath,
		input.Source,
		input.OAuthId,
		input.IDCardUrl,
		input.IDCardPath,
		input.ProvinceID,
		input.RegencyID,
		input.DistrictID,
		input.VillageID,
	)

	if err != nil {
		if f, ok := stacktrace.RootCause(err).(failure.Failure); ok {
			switch f.Code {
			case failure.CodeIncorrectUserID:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			}
		}

		log.WithContext(ctx).Error(stacktrace.Propagate(err, "auth update profile error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, "unknown error")
		return
	}

	respond.Success(c, trx, http.StatusCreated, input)
}

func (h *Handler) HandleChangePassword(c *gin.Context) {
	ctx := activity.NewContext("profile_change_password")
	ctx = activity.WithUserID(ctx, c.Value("user_id").(string))
	trx, _ := activity.GetTransactionID(ctx)
	userID, _ := activity.GetUserID(ctx)
	var input InputChangePassword

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

	id, err := uuid.Parse(userID)
	if err != nil {
		respond.Error(c, trx, http.StatusBadRequest, respond.ErrBadRequest, "invalid user id")
		return
	}

	if input.NewPassword != input.NewPasswordConfirmation {
		respond.Error(c, trx, http.StatusBadRequest, respond.ErrBadRequest, "new password confirmation not match")
		return
	}

	err = h.service.ChangePassword(
		id,
		input.OldPassword,
		input.NewPassword,
		input.NewPasswordConfirmation,
	)

	if err != nil {
		if f, ok := stacktrace.RootCause(err).(failure.Failure); ok {
			switch f.Code {
			case failure.CodeIncorrectUserID:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			case failure.CodeIncorrectPassword:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			}
		}

		log.WithContext(ctx).Error(stacktrace.Propagate(err, "auth update profile error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, "unknown error")
		return
	}

	respond.Success(c, trx, http.StatusCreated, nil)
}
