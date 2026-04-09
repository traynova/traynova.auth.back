package app

import (
	"context"
	"gestrym/src/common/models"
	"gestrym/src/core/roles/domain/ports"
	structs_roles "gestrym/src/core/roles/domain/structs"
)

type IRoleService interface {
	CreateRole(ctx context.Context, name, description string) (*structs_roles.RolesResponse, error)
	UpdateRole(ctx context.Context, id uint, name, description string) (*structs_roles.RolesResponse, error)
	DisableRole(ctx context.Context, id uint) error
	GetRoles() ([]structs_roles.RolesResponse, error)
}

type roleService struct {
	repo ports.IRoleRepository
}

func NewRoleService(r ports.IRoleRepository) IRoleService {
	return &roleService{repo: r}
}

func (s *roleService) CreateRole(ctx context.Context, name, description string) (*structs_roles.RolesResponse, error) {
	role := &models.Role{
		Name:        name,
		Description: description,
		IsActive:    true,
	}
	err := s.repo.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	rolesRespoonse := structs_roles.RolesResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
	return &rolesRespoonse, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id uint, name, description string) (*structs_roles.RolesResponse, error) {
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
	roleResponse := &structs_roles.RolesResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}

	return roleResponse, nil
}

func (s *roleService) DisableRole(ctx context.Context, id uint) error {
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	role.IsActive = false
	return s.repo.Disable(ctx, role)
}

func (s *roleService) GetRoles() ([]structs_roles.RolesResponse, error) {
	roles, err := s.repo.GetRoles()
	if err != nil {
		return nil, err
	}
	var rolesResponse []structs_roles.RolesResponse

	for _, role := range roles {
		rolesResponse = append(rolesResponse, structs_roles.RolesResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}
	return rolesResponse, nil
}
