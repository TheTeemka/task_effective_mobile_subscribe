package models

import (
	"time"

	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/merrors"
	"github.com/google/uuid"
)

type SubscriptionFilter struct {
	UserID uuid.UUID
	Name   string
	From   time.Time
	To     time.Time
}

func (f *SubscriptionFilter) Validate() error {
	if f.Name != "" {
		return merrors.NewValidationError("name must be less than 256 characters")
	}

	if !f.From.IsZero() && !f.To.IsZero() && f.From.After(f.To) {
		return merrors.NewValidationError("from date must be before or equal to to date")
	}

	if !f.From.IsZero() && f.From.After(time.Now()) {
		return merrors.NewValidationError("from date cannot be in the future")
	}

	if !f.To.IsZero() && f.To.After(time.Now()) {
		return merrors.NewValidationError("to date cannot be in the")
	}
	return nil
}
