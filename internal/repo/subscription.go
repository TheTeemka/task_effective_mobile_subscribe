package repo

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/models"
	"github.com/google/uuid"
)

type SubscriptionRepo struct {
	DB *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{DB: db}
}

func (s *SubscriptionRepo) Create(subscription *models.SubscriptionModel) error {
	query := `
        INSERT INTO subscriptions (user_id, cost, start_date, end_date, name)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := s.DB.Exec(query, subscription.UserID,
		subscription.Cost, subscription.StartDate, subscription.EndDate, subscription.Name)
	return err
}

func (s *SubscriptionRepo) FindByUserID(userID uuid.UUID) ([]*models.SubscriptionModel, error) {
	query := `
        SELECT id, user_id, cost, start_date, end_date, name
        FROM subscriptions
        WHERE user_id = $1`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.SubscriptionModel
	for rows.Next() {
		subscription := &models.SubscriptionModel{}
		err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.Cost,
			&subscription.StartDate, &subscription.EndDate, &subscription.Name)
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

func (s *SubscriptionRepo) FindByID(ID int64) (*models.SubscriptionModel, error) {
	subscription := &models.SubscriptionModel{}
	query := `
        SELECT id, user_id, cost, start_date, end_date, name
        FROM subscriptions
        WHERE id = $1`

	err := s.DB.QueryRow(query, ID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.Cost,
		&subscription.StartDate, &subscription.EndDate, &subscription.Name)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *SubscriptionRepo) Update(subscription *models.SubscriptionModel) error {
	query := `
        UPDATE subscriptions
        SET cost = $2, start_date = $3, end_date = $4, user_id = $5, name = $6
        WHERE id = $1`

	_, err := s.DB.Exec(query, subscription.ID, subscription.Cost, subscription.StartDate,
		subscription.EndDate, subscription.UserID, subscription.Name)
	return err
}

func (s *SubscriptionRepo) Delete(id int64) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

func (s *SubscriptionRepo) GetSum(filters *models.SubscriptionFilter) (float64, error) {
	var sum float64

	builder := squirrel.Select("COALESCE(SUM(cost), 0)").From("subscriptions")

	if filters.UserID != uuid.Nil {
		builder = builder.Where(squirrel.Eq{"user_id": filters.UserID})
	}

	if filters.Name != "" {
		builder = builder.Where(squirrel.Eq{"name": filters.Name})
	}

	if !filters.From.IsZero() {
		builder = builder.Where(squirrel.GtOrEq{"start_date": filters.From})
	}

	if !filters.To.IsZero() {
		builder = builder.Where(squirrel.LtOrEq{"end_date": filters.To})
	}

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
