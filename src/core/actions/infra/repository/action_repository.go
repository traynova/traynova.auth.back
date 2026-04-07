package repository

import (
	"context"
	"traynova/src/common/models"
	"traynova/src/core/actions/domain/ports"

	"gorm.io/gorm"
)

type actionRepository struct {
	db *gorm.DB
}

func NewActionRepository(db *gorm.DB) ports.IActionRepository {
	return &actionRepository{db: db}
}

func (r *actionRepository) Create(ctx context.Context, action *models.Action) error {
	return r.db.WithContext(ctx).Create(action).Error
}

func (r *actionRepository) FindAll(ctx context.Context) ([]models.Action, error) {
	var actions []models.Action
	err := r.db.WithContext(ctx).Find(&actions).Error
	return actions, err
}

func (r *actionRepository) FindByID(ctx context.Context, id uint) (*models.Action, error) {
	var action models.Action
	err := r.db.WithContext(ctx).First(&action, id).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}
