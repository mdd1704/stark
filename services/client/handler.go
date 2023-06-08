package client

import (
	"errors"
	"net/http"
	"strconv"

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

func (h *Handler) HandleCreate(c *gin.Context) {
	ctx := activity.NewContext("client_create")
	trx, _ := activity.GetTransactionID(ctx)
	var input Input

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

	user, err := h.service.Create(input.Name)
	if err != nil {
		if f, ok := stacktrace.RootCause(err).(failure.Failure); ok {
			switch f.Code {
			case failure.CodeUserAlreadyExist:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			}
		}

		log.WithContext(ctx).Error(stacktrace.Propagate(err, "create user error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, "unknown error")
		return
	}

	respond.Success(c, trx, http.StatusCreated, user)
}

func (h *Handler) HandleDetail(c *gin.Context) {
	ctx := activity.NewContext("client_detail")
	trx, _ := activity.GetTransactionID(ctx)
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, trx, http.StatusBadRequest, respond.ErrBadRequest, "invalid user id")
		return
	}

	user, err := h.service.FindByID(userID)
	if err != nil {
		if f, ok := stacktrace.RootCause(err).(failure.Failure); ok {
			switch f.Code {
			case failure.CodeUserNotFound:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			}
		}

		log.WithContext(ctx).Error(stacktrace.Propagate(err, "get user detail error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, err.Error())
		return
	}

	respond.Success(c, trx, http.StatusOK, user)
}

func (h *Handler) HandleUpdate(c *gin.Context) {
	ctx := activity.NewContext("client_update")
	trx, _ := activity.GetTransactionID(ctx)
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, trx, http.StatusBadRequest, respond.ErrBadRequest, "invalid user id")
		return
	}

	var input Input

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

	user, err := h.service.Update(userID, input.Name)
	if err != nil {
		if f, ok := stacktrace.RootCause(err).(failure.Failure); ok {
			switch f.Code {
			case failure.CodeUserNotFound:
				respond.Error(c, trx, http.StatusBadRequest, f.Code, f.Desc)
				return
			}
		}

		log.WithContext(ctx).Error(stacktrace.Propagate(err, "update user error"))
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, err.Error())
		return
	}

	respond.Success(c, trx, http.StatusCreated, user)
}

func (h *Handler) HandleAllByFilter(c *gin.Context) {
	ctx := activity.NewContext("client_all_by_filter")
	trx, _ := activity.GetTransactionID(ctx)
	var input Filter

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

	tenant, err := h.service.FindAllByFilter(input)
	if err != nil {
		respond.Error(c, trx, http.StatusInternalServerError, respond.ErrInternal, err.Error())
		return
	}

	respond.Success(c, trx, http.StatusCreated, tenant)
}

func (h *Handler) HandlePage(c *gin.Context) {
	ctx := activity.NewContext("client_page")
	trx, _ := activity.GetTransactionID(ctx)
	pageString := c.Query("page")
	limitString := c.Query("limit")
	page := 1
	limit := 25
	var err error
	if pageString != "" {
		page, err = strconv.Atoi(pageString)
		if err != nil {
			respond.Error(c, trx, http.StatusBadRequest, respond.ErrBadRequest, "invalid page")
			return
		}
	}

	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			respond.Error(c, trx, http.StatusBadRequest, respond.ErrBadRequest, "invalid limit")
			return
		}
	}

	tenantPage, err := h.service.FindPage(page, limit)
	if err != nil {
		respond.Error(c, trx, http.StatusBadRequest, respond.ErrInternal, err.Error())
		return
	}

	respond.Success(c, trx, http.StatusOK, tenantPage)
}
