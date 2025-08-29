package repo

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/merrors"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/models"
)

type SubscriptionRepo struct {
	DB *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{DB: db}
}

func (s *SubscriptionRepo) Create(subscription *models.SubscriptionModel) error {
	query := `
        INSERT INTO subscriptions (user_id, price, start_date, end_date, service_name)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := s.DB.Exec(query, subscription.UserID,
		subscription.Price, subscription.StartDate, subscription.EndDate, subscription.ServiceName)
	return err
}

func (s *SubscriptionRepo) GetByFilters(filters *models.SubscriptionFilter) ([]*models.SubscriptionModel, error) {

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id, user_id, price, start_date, end_date, service_name").
		From("subscriptions")

	builder = filters.ToSQL(builder)

	query, args, err := builder.ToSql()
	slog.Info("GetByFilters query", "query", query, "args", args)

	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.SubscriptionModel
	for rows.Next() {
		subscription := &models.SubscriptionModel{}
		err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.Price,
			&subscription.StartDate, &subscription.EndDate, &subscription.ServiceName)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionRepo) GetByID(ID int64) (*models.SubscriptionModel, error) {
	subscription := &models.SubscriptionModel{}
	query := `
        SELECT id, user_id, price, start_date, end_date, service_name
        FROM subscriptions
        WHERE id = $1`

	err := s.DB.QueryRow(query, ID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.Price,
		&subscription.StartDate, &subscription.EndDate, &subscription.ServiceName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, merrors.NewNotFoundErr("subscription not found")
		}
		return nil, err
	}
	return subscription, nil
}

func (s *SubscriptionRepo) Update(subscription *models.SubscriptionModel) error {
	query := `
        UPDATE subscriptions
        SET price = $2, start_date = $3, end_date = $4, user_id = $5, service_name = $6
        WHERE id = $1`

	res, err := s.DB.Exec(query, subscription.ID, subscription.Price,
		subscription.StartDate, subscription.EndDate, subscription.UserID,
		subscription.ServiceName)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return merrors.NewNotFoundErr("subscription not found")
	}
	return err
}

func (s *SubscriptionRepo) Delete(id int64) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	res, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return merrors.NewNotFoundErr("subscription not found")
	}
	return err
}

func (s *SubscriptionRepo) GetSum(filters *models.SubscriptionFilter) (float64, error) {
	var sum float64

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("COALESCE(SUM(price), 0)").
		From("subscriptions")

	builder = filters.ToSQL(builder)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	err = s.DB.QueryRow(query, args...).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}
