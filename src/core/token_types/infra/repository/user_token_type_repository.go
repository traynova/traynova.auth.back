package repository

import (
	"context"
	"gestrym/src/common/models"
	"gestrym/src/core/token_types/domain/ports"

	"gorm.io/gorm"
)

type userTokenTypeRepository struct {
	db *gorm.DB
}

func NewUserTokenTypeRepository(db *gorm.DB) ports.IUserTokenTypeRepository {
	return &userTokenTypeRepository{db: db}
}

func (r *userTokenTypeRepository) Create(ctx context.Context, tokenType *models.UserTokenType) error {
	return r.db.WithContext(ctx).Create(tokenType).Error
}

func (r *userTokenTypeRepository) FindAll(ctx context.Context) ([]models.UserTokenType, error) {
	var results []models.UserTokenType
	err := r.db.WithContext(ctx).Find(&results).Error
	return results, err
}

func (r *userTokenTypeRepository) FindByID(ctx context.Context, id uint) (*models.UserTokenType, error) {
	var result models.UserTokenType
	err := r.db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *userTokenTypeRepository) FindByType(ctx context.Context, typeName string) (*models.UserTokenType, error) {
	var result models.UserTokenType
	err := r.db.WithContext(ctx).Where("type = ?", typeName).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
