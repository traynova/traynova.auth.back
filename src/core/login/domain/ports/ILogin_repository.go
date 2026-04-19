package ports

import "gestrym/src/common/models"

type ILoginRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	UpdateInitialLogin(email string) error
}
