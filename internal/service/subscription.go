package service

import (
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/models"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/repo"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	*repo.SubscriptionRepo
}

func NewSubscriptionService(subscriptionRepo *repo.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{SubscriptionRepo: subscriptionRepo}
}

func (s *SubscriptionService) Create(subscription *models.SubscriptionCreateReq) error {
	if err := subscription.Validate(); err != nil {
		return err
	}

	subs := &models.SubscriptionModel{
		Name:      subscription.Name,
		UserID:    subscription.UserID,
		Cost:      subscription.Cost,
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
	}

	return s.SubscriptionRepo.Create(subs)
}

func (s *SubscriptionService) FindByUserID(uuid uuid.UUID) ([]*models.SubscriptionModel, error) {
	return s.SubscriptionRepo.FindByUserID(uuid)
}

func (s *SubscriptionService) FindByID(ID int64) (*models.SubscriptionModel, error) {
	return s.SubscriptionRepo.FindByID(ID)
}

func (s *SubscriptionService) Delete(ID int64) error {
	return s.SubscriptionRepo.Delete(ID)
}

func (s *SubscriptionService) GetSum(filter *models.SubscriptionFilter) (float64, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}
	return s.SubscriptionRepo.GetSum(filter)
}
