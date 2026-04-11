package ports_auth

import "gestrym/src/common/models"

type IAuthRepository interface {
	ValidateEmail(email string) (*models.User, error)
	UpdateUSer(user *models.User) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	ValidateCoachGymAssociation(coachId uint) (*models.TrainerProfile, error)
	GetAllUsers(page int, pageSize int, name *string, dni *string, email *string) ([]models.User, int64, error)
}
