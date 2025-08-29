package models

import (
	"time"

	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/merrors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionModel struct {
	ID        int64
	Name      string `json:"name" validate:"required"`
	UserID    uuid.UUID
	Cost      float64
	StartDate time.Time
	EndDate   time.Time
}

type SubscriptionCreateReq struct {
	Name      string    `json:"name" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Cost      float64   `json:"cost" validate:"required,gt=0"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (s *SubscriptionCreateReq) Validate() error {
	if err := validate.Struct(s); err != nil {
		return err
	}

	if s.UserID == uuid.Nil {
		return merrors.NewValidationError("user_id cannot be nil")
	}

	if !s.EndDate.IsZero() && s.EndDate.Before(s.StartDate) {
		return merrors.NewValidationError("end_date must be after start_date")
	}

	return nil
}
