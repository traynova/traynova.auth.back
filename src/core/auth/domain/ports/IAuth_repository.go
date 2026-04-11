package ports_auth

import "gestrym/src/common/models"

type IAuthRepository interface {
	ValidateEmail(email string) (*models.User, error)
	UpdateUSer(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
	ValidateCoachGymAssociation(coachId uint) (*models.TrainerProfile, error)
	CreateTrainerProfile(profile *models.TrainerProfile) (*models.TrainerProfile, error)
	UpdateTrainerProfile(profile *models.TrainerProfile) (*models.TrainerProfile, error)
	GetTrainerProfilesByUserID(userID uint) ([]models.TrainerProfile, error)
	GetTrainerProfileByUserIDAndGymID(userID uint, gymID *uint) (*models.TrainerProfile, error)
	CreateGymProfile(profile *models.GymProfile) (*models.GymProfile, error)
	UpdateGymProfile(profile *models.GymProfile) (*models.GymProfile, error)
	GetGymProfileByUserID(userID uint) (*models.GymProfile, error)
	CreateTrainerClient(client *models.TrainerClient) (*models.TrainerClient, error)
	GetTrainerClientByProfileAndClient(profileID uint, clientID uint) (*models.TrainerClient, error)
	GetTrainerClientsByProfileID(profileID uint) ([]models.User, error)
	GetGymTrainersByGymUserID(gymUserID uint) ([]models.TrainerProfile, error)
	CreateGymClient(client *models.GymClient) (*models.GymClient, error)
	GetGymClientByGymAndClient(gymUserID uint, clientID uint) (*models.GymClient, error)
	GetAllUsers(page int, pageSize int, name *string, dni *string, email *string) ([]models.User, int64, error)
}
