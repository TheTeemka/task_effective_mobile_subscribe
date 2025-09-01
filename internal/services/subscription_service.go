package services

import (
	"fmt"

	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/models"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/repo"
)

type SubscriptionService struct {
	subscriptionRepo *repo.SubscriptionRepo
}

func NewSubscriptionService(subscriptionRepo *repo.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{subscriptionRepo: subscriptionRepo}
}

func (s *SubscriptionService) Create(subCreateReq *models.SubscriptionCreateReq) error {
	if err := subCreateReq.Validate(); err != nil {
		return fmt.Errorf("subscription creation validation failed: %w", err)
	}

	sub, err := subCreateReq.ToModel()
	if err != nil {
		return fmt.Errorf("failed to convert subscription request to model: %w", err)
	}

	return s.subscriptionRepo.Create(sub)
}

func (s *SubscriptionService) GetSum(filter *models.SubscriptionFilter) (float64, error) {
	if err := filter.Validate(); err != nil {
		return 0, fmt.Errorf("subscription sum filter validation failed: %w", err)
	}

	sum, err := s.subscriptionRepo.GetSum(filter)
	if err != nil {
		return 0, fmt.Errorf("failed to get subscription sum from repository: %w", err)
	}

	return sum, nil
}

func (s *SubscriptionService) GetByFilters(filter *models.SubscriptionFilter) ([]*models.SubscriptionModel, error) {
	subs, err := s.subscriptionRepo.GetByFilters(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by filters: %w", err)
	}

	return subs, nil
}

func (s *SubscriptionService) GetByID(ID int64) (*models.SubscriptionModel, error) {
	sub, err := s.subscriptionRepo.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription by ID %d: %w", ID, err)
	}

	return sub, nil
}

func (s *SubscriptionService) Delete(ID int64) error {
	err := s.subscriptionRepo.Delete(ID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription with ID %d: %w", ID, err)
	}

	return nil
}

func (s *SubscriptionService) Update(id int64, subUpdateReq *models.SubscriptionUpdateReq) error {
	if err := subUpdateReq.Validate(); err != nil {
		return fmt.Errorf("subscription update validation failed: %w", err)
	}

	sub, err := s.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get existing subscription for update: %w", err)
	}

	err = subUpdateReq.PatchModel(sub)
	if err != nil {
		return fmt.Errorf("failed to patch subscription model: %w", err)
	}

	return s.subscriptionRepo.Update(sub)
}
