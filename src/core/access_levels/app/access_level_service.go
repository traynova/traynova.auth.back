package app

import (
	"context"
	"traynova/src/common/models"
	"traynova/src/core/access_levels/domain/ports"
)

type IAccessLevelService interface {
	CreateAccessLevel(ctx context.Context, name, description string) (*models.AccessLevel, error)
	GetAccessLevels(ctx context.Context) ([]models.AccessLevel, error)
}

type accessLevelService struct {
	repo ports.IAccessLevelRepository
}

func NewAccessLevelService(r ports.IAccessLevelRepository) IAccessLevelService {
	return &accessLevelService{repo: r}
}

func (s *accessLevelService) CreateAccessLevel(ctx context.Context, name, description string) (*models.AccessLevel, error) {
	level := &models.AccessLevel{Name: name, Description: description}
	err := s.repo.Create(ctx, level)
	if err != nil {
		return nil, err
	}
	return level, nil
}

func (s *accessLevelService) GetAccessLevels(ctx context.Context) ([]models.AccessLevel, error) {
	return s.repo.FindAll(ctx)
}
