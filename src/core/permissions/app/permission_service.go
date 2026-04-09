package app

import (
	"context"
	"gestrym/src/common/models"
	"gestrym/src/core/permissions/domain/ports"
)

type IPermissionService interface {
	CreatePermission(ctx context.Context, roleID, actionID, resourceID uint) (*models.Permission, error)
}

type permissionService struct {
	repo ports.IPermissionRepository
}

func NewPermissionService(r ports.IPermissionRepository) IPermissionService {
	return &permissionService{repo: r}
}

func (s *permissionService) CreatePermission(ctx context.Context, roleID, actionID, resourceID uint) (*models.Permission, error) {
	permission := &models.Permission{
		RoleID:     roleID,
		ActionID:   actionID,
		ResourceID: resourceID,
		IsActive:   true,
	}
	err := s.repo.Create(ctx, permission)
	if err != nil {
		return nil, err
	}
	return permission, nil
}
