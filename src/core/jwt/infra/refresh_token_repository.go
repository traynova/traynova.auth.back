package jwt_repository

import (
	"gestrym/src/common/models"
	jwt_ports "gestrym/src/core/jwt/domain/ports"
	"sync"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

var (
	refreshTokenRepositoryInstance *RefreshTokenRepository
	refreshTokenRepositoryOnce     sync.Once
)

func NewRefreshTokenRepository(db *gorm.DB) jwt_ports.IRefreshTokenRepository {
	refreshTokenRepositoryOnce.Do(func() {
		refreshTokenRepositoryInstance = &RefreshTokenRepository{
			db: db,
		}
	})
	return refreshTokenRepositoryInstance
}

func (r *RefreshTokenRepository) InsertRefreshToken(refreshToken *models.RefreshToken) error {
	err := r.db.Create(&refreshToken).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) UpdateRefreshToken(refreshToken *models.RefreshToken) error {
	err := r.db.Save(&refreshToken).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByKey(key string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("key = ?", key).First(&refreshToken).Error
	if err != nil {
		return refreshToken, err
	}
	return refreshToken, nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByUserId(userId uint) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("user_id = ?", userId).First(&refreshToken).Error
	if err != nil {
		return refreshToken, err
	}
	return refreshToken, nil
}

func (r *RefreshTokenRepository) GetRefreshTokenWithUserAndRole(key string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("key = ?", key).Preload("User.Role").First(&refreshToken).Error
	if err != nil {
		return refreshToken, err
	}
	return refreshToken, nil
}
