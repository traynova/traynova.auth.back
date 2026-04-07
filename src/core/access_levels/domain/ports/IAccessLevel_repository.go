package ports

import (
	"context"
	"traynova/src/common/models"
)

type IAccessLevelRepository interface {
	Create(ctx context.Context, level *models.AccessLevel) error
	FindAll(ctx context.Context) ([]models.AccessLevel, error)
	FindByID(ctx context.Context, id uint) (*models.AccessLevel, error)
}
