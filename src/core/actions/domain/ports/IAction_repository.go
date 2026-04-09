package ports

import (
	"context"
	"gestrym/src/common/models"
)

type IActionRepository interface {
	Create(ctx context.Context, action *models.Action) error
	FindAll(ctx context.Context) ([]models.Action, error)
	FindByID(ctx context.Context, id uint) (*models.Action, error)
}
