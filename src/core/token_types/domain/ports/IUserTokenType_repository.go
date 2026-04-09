package ports

import (
	"context"
	"gestrym/src/common/models"
)

type IUserTokenTypeRepository interface {
	Create(ctx context.Context, tokenType *models.UserTokenType) error
	FindAll(ctx context.Context) ([]models.UserTokenType, error)
	FindByID(ctx context.Context, id uint) (*models.UserTokenType, error)
	FindByType(ctx context.Context, t string) (*models.UserTokenType, error)
}
