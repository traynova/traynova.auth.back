package app

import (
	"context"
	"gestrym/src/common/models"
	"gestrym/src/core/actions/domain/ports"
)

type IActionService interface {
	CreateAction(ctx context.Context, name string) (*models.Action, error)
	GetActions(ctx context.Context) ([]models.Action, error)
}

type actionService struct {
	repo ports.IActionRepository
}

func NewActionService(r ports.IActionRepository) IActionService {
	return &actionService{repo: r}
}

func (s *actionService) CreateAction(ctx context.Context, name string) (*models.Action, error) {
	action := &models.Action{Name: name}
	err := s.repo.Create(ctx, action)
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (s *actionService) GetActions(ctx context.Context) ([]models.Action, error) {
	return s.repo.FindAll(ctx)
}
