package models

import (
	"net/url"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/merrors"
	"github.com/google/uuid"
)

type SubscriptionFilter struct {
	UserID      uuid.UUID
	ServiceName string
	From        time.Time
	To          time.Time
}

func NewSubscriptionFilterFromURL(q url.Values) (*SubscriptionFilter, error) {
	filter := &SubscriptionFilter{}

	if userIDStr := q.Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, merrors.NewValidationError("Invalid user_id")
		}
		filter.UserID = userID
	}

	if name := q.Get("service_name"); name != "" {
		filter.ServiceName = name
	}

	if fromStr := q.Get("from"); fromStr != "" {
		from, err := time.Parse(TimeFormat, fromStr)
		if err != nil {
			return nil, merrors.NewValidationError("Invalid from date")
		}
		filter.From = from
	}

	if toStr := q.Get("to"); toStr != "" {
		to, err := time.Parse(TimeFormat, toStr)
		if err != nil {
			return nil, merrors.NewValidationError("Invalid to date")
		}
		filter.To = to
	}

	return filter, nil
}

func (f *SubscriptionFilter) Validate() error {
	if !f.From.IsZero() && !f.To.IsZero() && f.From.After(f.To) {
		return merrors.NewValidationError("from date must be before or equal to to date")
	}
	return nil
}

func (f *SubscriptionFilter) ToSQL(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if f.UserID != uuid.Nil {
		builder = builder.Where(squirrel.Eq{"user_id": f.UserID})
	}

	if f.ServiceName != "" {
		builder = builder.Where(squirrel.Eq{"service_name": f.ServiceName})
	}

	if !f.From.IsZero() {
		builder = builder.Where(squirrel.GtOrEq{"start_date": f.From})
	}

	if !f.To.IsZero() {
		builder = builder.Where(squirrel.LtOrEq{"end_date": f.To})
	}

	return builder
}
