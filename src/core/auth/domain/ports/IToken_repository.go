package ports

import (
	"context"
	"gestrym/src/common/models"
)

type ITokenRepository interface {
	CreateUserToken(ctx context.Context, ut *models.UserToken) error
	CreateRefreshToken(ctx context.Context, rt *models.RefreshToken) error
	FindUserToken(ctx context.Context, token string) (*models.UserToken, error)
	FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeUserToken(ctx context.Context, token string) error
	RevokeRefreshToken(ctx context.Context, token string) error
}
