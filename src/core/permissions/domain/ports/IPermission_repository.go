package ports

import (
	"context"
	"gestrym/src/common/models"
)

type IPermissionRepository interface {
	Create(ctx context.Context, permission *models.Permission) error
	FindAll(ctx context.Context) ([]models.Permission, error)
}
