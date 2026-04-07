package app

import (
	"context"
	"traynova/src/common/models"
	"traynova/src/core/token_types/domain/ports"
)

type IUserTokenTypeService interface {
	CreateTokenType(ctx context.Context, typeName string) (*models.UserTokenType, error)
	GetTokenTypes(ctx context.Context) ([]models.UserTokenType, error)
}

type userTokenTypeService struct {
	repo ports.IUserTokenTypeRepository
}

func NewUserTokenTypeService(r ports.IUserTokenTypeRepository) IUserTokenTypeService {
	return &userTokenTypeService{repo: r}
}

func (s *userTokenTypeService) CreateTokenType(ctx context.Context, typeName string) (*models.UserTokenType, error) {
	tt := &models.UserTokenType{Type: typeName}
	err := s.repo.Create(ctx, tt)
	if err != nil {
		return nil, err
	}
	return tt, nil
}

func (s *userTokenTypeService) GetTokenTypes(ctx context.Context) ([]models.UserTokenType, error) {
	return s.repo.FindAll(ctx)
}
