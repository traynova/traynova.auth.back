package repository

import (
	"context"
	"traynova/src/common/models"
	"traynova/src/core/auth/domain/ports"

	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) ports.ITokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) CreateUserToken(ctx context.Context, ut *models.UserToken) error {
	return r.db.WithContext(ctx).Create(ut).Error
}

func (r *tokenRepository) CreateRefreshToken(ctx context.Context, rt *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(rt).Error
}

func (r *tokenRepository) FindUserToken(ctx context.Context, token string) (*models.UserToken, error) {
	var ut models.UserToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&ut).Error
	if err != nil {
		return nil, err
	}
	return &ut, nil
}

func (r *tokenRepository) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *tokenRepository) RevokeUserToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&models.UserToken{}).Where("token = ?", token).Update("is_revoked", true).Error
}

func (r *tokenRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).Where("token = ?", token).Update("is_revoked", true).Error
}
