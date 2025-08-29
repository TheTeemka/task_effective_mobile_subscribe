package models

import (
	"time"

	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/merrors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const TimeFormat = "01-2006"

type SubscriptionModel struct {
	ID          int64
	ServiceName string
	UserID      uuid.UUID
	Price       float64
	StartDate   time.Time
	EndDate     *time.Time
}

type SubscriptionCreateReq struct {
	ServiceName string    `json:"service_name" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
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

	startDate, err := time.Parse(TimeFormat, s.StartDate)
	if err != nil {
		return merrors.NewValidationError("invalid start_date format (expected MM-YYYY)")
	}

	if s.EndDate != nil {
		endDate, err := time.Parse(TimeFormat, *s.EndDate)
		if err != nil {
			return merrors.NewValidationError("invalid end_date format (expected MM-YYYY)")
		}
		if endDate.Before(startDate) {
			return merrors.NewValidationError("end_date must be after start_date")
		}
	}

	return nil
}

func (s *SubscriptionCreateReq) ToModel() (*SubscriptionModel, error) {
	subModel := &SubscriptionModel{
		ServiceName: s.ServiceName,
		UserID:      s.UserID,
		Price:       s.Price,
	}

	var err error
	if s.StartDate != "" {
		subModel.StartDate, err = time.Parse(TimeFormat, s.StartDate)
		if err != nil {
			return nil, merrors.NewValidationError("invalid start_date format (expected MM-YYYY)")
		}
	}

	if s.EndDate != nil {
		endDate, err := time.Parse(TimeFormat, *s.EndDate)
		if err != nil {
			return nil, merrors.NewValidationError("invalid start_date format (expected MM-YYYY)")
		}
		subModel.EndDate = &endDate
	}

	return subModel, nil
}

type SubscriptionUpdateReq struct {
	ServiceName *string    `json:"service_name"`
	UserID      *uuid.UUID `json:"user_id"`
	Price       *float64   `json:"price"`
	StartDate   *string    `json:"start_date"`
	EndDate     *string    `json:"end_date"`
}

func (s *SubscriptionUpdateReq) Validate() error {
	if s.UserID != nil && *s.UserID == uuid.Nil {
		return merrors.NewValidationError("user_id cannot be nil")
	}

	if s.StartDate != nil && s.EndDate != nil {
		startDate, err := time.Parse(TimeFormat, *s.StartDate)
		if err != nil {
			return merrors.NewValidationError("invalid start_date format (expected MM-YYYY)")
		}
		endDate, err := time.Parse(TimeFormat, *s.EndDate)
		if err != nil {
			return merrors.NewValidationError("invalid end_date format (expected MM-YYYY)")
		}
		if endDate.Before(startDate) {
			return merrors.NewValidationError("end_date must be after start_date")
		}
	}

	return nil
}

func (s *SubscriptionUpdateReq) PatchModel(subModel *SubscriptionModel) error {
	var err error
	if s.ServiceName != nil {
		subModel.ServiceName = *s.ServiceName
	}
	if s.UserID != nil {
		subModel.UserID = *s.UserID
	}
	if s.Price != nil {
		subModel.Price = *s.Price
	}
	if s.StartDate != nil {
		subModel.StartDate, err = time.Parse(TimeFormat, *s.StartDate)
		if err != nil {
			return merrors.NewValidationError("invalid start_date format (expected MM-YYYY)")
		}
	}
	if s.EndDate != nil {
		endDate, err := time.Parse(TimeFormat, *s.EndDate)
		if err != nil {
			return merrors.NewValidationError("invalid start_date format (expected MM-YYYY)")
		}
		subModel.EndDate = &endDate
	}
	return nil
}
