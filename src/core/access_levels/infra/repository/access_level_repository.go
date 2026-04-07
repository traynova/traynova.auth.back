package repository

import (
	"context"
	"traynova/src/common/models"
	"traynova/src/core/access_levels/domain/ports"

	"gorm.io/gorm"
)

type accessLevelRepository struct {
	db *gorm.DB
}

func NewAccessLevelRepository(db *gorm.DB) ports.IAccessLevelRepository {
	return &accessLevelRepository{db: db}
}

func (r *accessLevelRepository) Create(ctx context.Context, action *models.AccessLevel) error {
	return r.db.WithContext(ctx).Create(action).Error
}

func (r *accessLevelRepository) FindAll(ctx context.Context) ([]models.AccessLevel, error) {
	var actions []models.AccessLevel
	err := r.db.WithContext(ctx).Find(&actions).Error
	return actions, err
}

func (r *accessLevelRepository) FindByID(ctx context.Context, id uint) (*models.AccessLevel, error) {
	var action models.AccessLevel
	err := r.db.WithContext(ctx).First(&action, id).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}
