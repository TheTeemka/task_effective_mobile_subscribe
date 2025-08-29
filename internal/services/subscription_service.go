package services

import (
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
		return err
	}

	sub, err := subCreateReq.ToModel()
	if err != nil {
		return err
	}

	return s.subscriptionRepo.Create(sub)
}

func (s *SubscriptionService) GetSum(filter *models.SubscriptionFilter) (float64, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}
	return s.subscriptionRepo.GetSum(filter)
}

func (s *SubscriptionService) GetByFilters(filter *models.SubscriptionFilter) ([]*models.SubscriptionModel, error) {
	return s.subscriptionRepo.GetByFilters(filter)
}

func (s *SubscriptionService) GetByID(ID int64) (*models.SubscriptionModel, error) {
	return s.subscriptionRepo.GetByID(ID)
}

func (s *SubscriptionService) Delete(ID int64) error {
	return s.subscriptionRepo.Delete(ID)
}

func (s *SubscriptionService) Update(id int64, subUpdateReq *models.SubscriptionUpdateReq) error {
	if err := subUpdateReq.Validate(); err != nil {
		return err
	}

	sub, err := s.GetByID(id)
	if err != nil {
		return err
	}

	err = subUpdateReq.PatchModel(sub)
	if err != nil {
		return err
	}
	return s.subscriptionRepo.Update(sub)
}
