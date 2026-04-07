package ports

import (
	"context"
	"traynova/src/common/models"
)

type IRoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	FindByID(ctx context.Context, id uint) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindAll(ctx context.Context) ([]models.Role, error)
	Update(ctx context.Context, role *models.Role) error
	Disable(ctx context.Context, role *models.Role) error
}
