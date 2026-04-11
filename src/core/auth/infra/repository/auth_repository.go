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
