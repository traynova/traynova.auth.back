package repository

import (
	"gestrym/src/common/models"
	ports_auth "gestrym/src/core/auth/domain/ports"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) ports_auth.IAuthRepository {
	return &authRepository{db}
}

func (a *authRepository) ValidateEmail(email string) (*models.User, error) {
	var user models.User
	err := a.db.Model(&models.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepository) UpdateUSer(user *models.User) (*models.User, error) {
	if err := a.db.Model(&models.User{}).Where("id = ?", user.ID).Unscoped().UpdateColumn("IsActive", true).Error; err != nil {
		return nil, err
	}

	if err := a.db.Model(&user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (a *authRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := a.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authRepository) UpdateUser(user *models.User) (*models.User, error) {
	if err := a.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := a.db.Preload("Role").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepository) ValidateCoachGymAssociation(coachId uint) (*models.TrainerProfile, error) {
	var association models.TrainerProfile
	err := a.db.Where("user_id = ?", coachId).First(&association).Error
	if err != nil {
		return nil, err
	}
	return &association, nil
}

func (a *authRepository) GetAllUsers(page int, pageSize int, name *string, dni *string, email *string) ([]models.User, int64, error) {
	offset := (page - 1) * pageSize

	query := a.db.
		Preload("Role").
		Where("is_active = ?", true).
		Order("id asc")

	// Filtrar por nombre
	if name != nil && *name != "" {
		query = query.Where("full_name ILIKE ?", "%"+*name+"%")
	}

	// Filtrar por DNI
	if dni != nil && *dni != "" {
		query = query.Where("dni ILIKE ?", "%"+*dni+"%")
	}

	// Filtrar por email
	if email != nil && *email != "" {
		query = query.Where("email ILIKE ?", "%"+*email+"%")
	}

	var total int64
	if err := query.Model(&models.User{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Offset(offset).Limit(pageSize)

	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (a *authRepository) CreateTrainerProfile(profile *models.TrainerProfile) (*models.TrainerProfile, error) {
	if err := a.db.Create(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (a *authRepository) UpdateTrainerProfile(profile *models.TrainerProfile) (*models.TrainerProfile, error) {
	if err := a.db.Save(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (a *authRepository) CreateGymProfile(profile *models.GymProfile) (*models.GymProfile, error) {
	if err := a.db.Create(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (a *authRepository) UpdateGymProfile(profile *models.GymProfile) (*models.GymProfile, error) {
	if err := a.db.Save(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (a *authRepository) GetGymProfileByUserID(userID uint) (*models.GymProfile, error) {
	var profile models.GymProfile
	if err := a.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (a *authRepository) GetTrainerProfilesByUserID(userID uint) ([]models.TrainerProfile, error) {
	var profiles []models.TrainerProfile
	if err := a.db.Preload("User").Where("user_id = ?", userID).Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

func (a *authRepository) GetTrainerProfileByUserIDAndGymID(userID uint, gymID *uint) (*models.TrainerProfile, error) {
	var profile models.TrainerProfile
	query := a.db.Preload("User").Where("user_id = ?", userID)
	if gymID == nil {
		query = query.Where("gim_id IS NULL")
	} else {
		query = query.Where("gim_id = ?", *gymID)
	}

	if err := query.First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (a *authRepository) CreateTrainerClient(client *models.TrainerClient) (*models.TrainerClient, error) {
	if err := a.db.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (a *authRepository) GetTrainerClientByProfileAndClient(profileID uint, clientID uint) (*models.TrainerClient, error) {
	var trainerClient models.TrainerClient
	if err := a.db.Where("trainer_profile_id = ? AND client_id = ?", profileID, clientID).First(&trainerClient).Error; err != nil {
		return nil, err
	}
	return &trainerClient, nil
}

func (a *authRepository) GetTrainerClientsByProfileID(profileID uint) ([]models.User, error) {
	var users []models.User
	err := a.db.Model(&models.User{}).
		Select("users.*").
		Joins("JOIN trainer_clients ON trainer_clients.client_id = users.id").
		Where("trainer_clients.trainer_profile_id = ?", profileID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (a *authRepository) GetGymTrainersByGymUserID(gymUserID uint) ([]models.TrainerProfile, error) {
	var profiles []models.TrainerProfile
	err := a.db.Preload("User").Where("gim_id = ?", gymUserID).Find(&profiles).Error
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (a *authRepository) CreateGymClient(client *models.GymClient) (*models.GymClient, error) {
	if err := a.db.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (a *authRepository) GetGymClientByGymAndClient(gymUserID uint, clientID uint) (*models.GymClient, error) {
	var gymClient models.GymClient
	if err := a.db.Where("gym_user_id = ? AND client_id = ?", gymUserID, clientID).First(&gymClient).Error; err != nil {
		return nil, err
	}
	return &gymClient, nil
}
