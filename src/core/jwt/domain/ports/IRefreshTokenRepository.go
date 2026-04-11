package jwt_ports

import "gestrym/src/common/models"

type IRefreshTokenRepository interface {
	GetRefreshTokenByKey(key string) (models.RefreshToken, error)
	GetRefreshTokenByUserId(userId uint) (models.RefreshToken, error)
	InsertRefreshToken(refreshToken *models.RefreshToken) error
	UpdateRefreshToken(refreshToken *models.RefreshToken) error
	GetRefreshTokenWithUserAndRole(key string) (models.RefreshToken, error)
}
