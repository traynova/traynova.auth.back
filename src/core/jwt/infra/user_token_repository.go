package jwt_repository

import (
	"gestrym/src/common/models"
	jwt_ports "gestrym/src/core/jwt/domain/ports"
	"sync"

	"gorm.io/gorm"
)

type userTokenRepository struct {
	db *gorm.DB
}

var (
	userTokenRepositoryInstance *userTokenRepository
	userTokenRepositoryOnce     sync.Once
)

func NewUserTokenRepository(db *gorm.DB) jwt_ports.IUserTokenRepository {
	userTokenRepositoryOnce.Do(func() {
		userTokenRepositoryInstance = &userTokenRepository{}
		userTokenRepositoryInstance.db = db
	})
	return userTokenRepositoryInstance
}

func (u *userTokenRepository) InsertUserToken(token models.UserToken) error {
	err := u.db.Create(&token).Error
	return err
}

func (u *userTokenRepository) DeleteUserToken(token string) error {
	userToken := &models.UserToken{}
	err := u.db.Unscoped().Where("token = ?", token).Delete(&userToken).Error
	return err
}

func (u *userTokenRepository) GetUserToken(token string) (models.UserToken, error) {
	userToken := &models.UserToken{}
	if err := u.db.Where("token = ?", token).First(&userToken).Error; err != nil {
		return models.UserToken{}, err
	}
	return *userToken, nil
}

func (r *userTokenRepository) GetLastActivationUserToken(userId uint) (string, error) {
	var validationToken models.UserToken
	err := r.db.Model(&models.UserToken{}).
		Joins("JOIN user_token_types ON user_token_types.id = user_tokens.user_token_type_id").
		Where("user_tokens.user_id = ? AND user_token_types.type = ?", userId, models.UserTokenTypeActivation).
		Last(&validationToken).Error
	if err != nil {
		return "", err
	}
	return validationToken.Token, nil
}
