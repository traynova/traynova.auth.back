package jwt_ports

import "gestrym/src/common/models"

type IUserTokenRepository interface {
	InsertUserToken(token models.UserToken) error
	DeleteUserToken(token string) error
	GetUserToken(token string) (models.UserToken, error)
	GetLastActivationUserToken(userId uint) (string, error)
}
