package app

import (
	"context"
	"traynova/src/common/models"
	"traynova/src/core/roles/domain/ports"
)

type IRoleService interface {
	CreateRole(ctx context.Context, name, description string) (*models.Role, error)
	UpdateRole(ctx context.Context, id uint, name, description string) (*models.Role, error)
	DisableRole(ctx context.Context, id uint) error
}

type roleService struct {
	repo ports.IRoleRepository
}

func NewRoleService(r ports.IRoleRepository) IRoleService {
	return &roleService{repo: r}
}

func (s *roleService) CreateRole(ctx context.Context, name, description string) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		Description: description,
		IsActive:    true,
	}
	err := s.repo.Create(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id uint, name, description string) (*models.Role, error) {
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		role.Name = name
	}
	if description != "" {
		role.Description = description
	}

	err = s.repo.Update(ctx, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) DisableRole(ctx context.Context, id uint) error {
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	role.IsActive = false
	return s.repo.Disable(ctx, role)
}
